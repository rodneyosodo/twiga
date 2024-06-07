// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/rodneyosodo/twiga/internal/jaeger"
	pgclient "github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/internal/prometheus"
	"github.com/rodneyosodo/twiga/internal/server"
	grpcserver "github.com/rodneyosodo/twiga/internal/server/grpc"
	httpserver "github.com/rodneyosodo/twiga/internal/server/http"
	"github.com/rodneyosodo/twiga/users"
	grpcapi "github.com/rodneyosodo/twiga/users/api/grpc"
	httpapi "github.com/rodneyosodo/twiga/users/api/http"
	"github.com/rodneyosodo/twiga/users/jwt"
	"github.com/rodneyosodo/twiga/users/middleware"
	"github.com/rodneyosodo/twiga/users/proto"
	"github.com/rodneyosodo/twiga/users/repository"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	svcName       = "users"
	defHTTPPort   = "6000"
	defGRPCPort   = "6010"
	envPrefixHTTP = "TWIGA_HTTP_"
	envPrefixGRPC = "TWIGA_GRPC_"
	envPrefixDB   = "TWIGA_DB_"
	defDB         = "users"
)

type config struct {
	LogLevel   string        `env:"TWIGA_LOG_LEVEL"          envDefault:"info"`
	JWTSecret  string        `env:"TWIGA_JWT_SECRET"         envDefault:"secret"`
	JWTExp     time.Duration `env:"TWIGA_JWT_EXP"            envDefault:"24h"`
	JaegerURL  url.URL       `env:"TWIGA_JAEGER_URL"         envDefault:"http://localhost:14268"`
	TraceRatio float64       `env:"TWIGA_JAEGER_TRACE_RATIO" envDefault:"1.0"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err.Error())
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		log.Fatalf("failed to parse log level: %s", err.Error())
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(logHandler)

	dbConfig := pgclient.Config{Name: defDB}
	if err := env.ParseWithOptions(&dbConfig, env.Options{Prefix: envPrefixDB}); err != nil {
		logger.Error(err.Error())
	}

	db, err := pgclient.Setup(dbConfig, *repository.Migration())
	if err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}
	defer db.Close()

	tp, err := jaeger.NewProvider(ctx, svcName, cfg.JaegerURL, cfg.TraceRatio)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init Jaeger: %s", err))
		cancel()
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)

	svc := newService(db, tracer, cfg, logger)

	httpServerConfig := server.Config{Port: defHTTPPort}
	if err := env.ParseWithOptions(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err.Error()))
		cancel()
		os.Exit(1)
	}
	hs := httpserver.New(ctx, cancel, svcName, httpServerConfig, httpapi.MakeHandler(svc), logger)

	grpcServerConfig := server.Config{Port: defGRPCPort}
	if err := env.ParseWithOptions(&grpcServerConfig, env.Options{Prefix: envPrefixGRPC}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s gRPC server configuration : %s", svcName, err.Error()))
		cancel()
		os.Exit(1)
	}
	registerAuthServiceServer := func(srv *grpc.Server) {
		reflection.Register(srv)
		proto.RegisterUsersServiceServer(srv, grpcapi.NewServer(svc))
	}

	gs := grpcserver.New(ctx, cancel, svcName, grpcServerConfig, registerAuthServiceServer, logger)

	g.Go(func() error {
		return hs.Start()
	})
	g.Go(func() error {
		return gs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs, gs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("%s service terminated: %s", svcName, err))
	}
}

func newService(db pgclient.Database, tracer trace.Tracer, cfg config, logger *slog.Logger) users.Service {
	urepo := repository.NewUsersRepository(db)
	prepo := repository.NewPreferencesRepository(db)
	forepo := repository.NewFollowingRepository(db)
	frepo := repository.NewFeedRepository(db)
	tokenizer := jwt.NewTokenizer(cfg.JWTSecret, cfg.JWTExp)

	svc := users.NewService(urepo, prepo, forepo, frepo, tokenizer)
	svc = middleware.NewLoggingMiddleware(logger, svc)
	svc = middleware.NewTracingMiddleware(tracer, svc)
	counter, latency := prometheus.MakeMetrics("users", "api")
	svc = middleware.NewMetricsMiddleware(counter, latency, svc)

	return svc
}
