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

	"github.com/caarlos0/env/v10"
	"github.com/chenjiandongx/ginprom"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	"github.com/rodneyosodo/twiga/internal/auth"
	"github.com/rodneyosodo/twiga/internal/jaeger"
	"github.com/rodneyosodo/twiga/internal/server"
	httpserver "github.com/rodneyosodo/twiga/internal/server/http"
	"github.com/rodneyosodo/twiga/posts"
	"github.com/rodneyosodo/twiga/posts/api"
	"github.com/rodneyosodo/twiga/posts/repository"
	"github.com/rodneyosodo/twiga/users/proto"
	sloggin "github.com/samber/slog-gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"golang.org/x/sync/errgroup"
)

const (
	svcName       = "posts"
	defHTTPPort   = "6001"
	envPrefixHTTP = "TWIGA_POSTS_HTTP_"
	envPrefixAuth = "TWIGA_USERS_GRPC_"
)

type config struct {
	LogLevel   string  `env:"TWIGA_POSTS_LOG_LEVEL"    envDefault:"info"`
	MongoURL   string  `env:"TWIGA_POSTS_MONGO_URL"    envDefault:"mongodb://localhost:27017"`
	MongoDB    string  `env:"TWIGA_POSTS_MONGO_DB"     envDefault:"twiga"`
	MongoColl  string  `env:"TWIGA_POSTS_MONGO_COLL"   envDefault:"posts"`
	JaegerURL  url.URL `env:"TWIGA_JAEGER_URL"         envDefault:"http://localhost:14268"`
	TraceRatio float64 `env:"TWIGA_JAEGER_TRACE_RATIO" envDefault:"1.0"`
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

	collection, err := connectDB(ctx, cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to MongoDB: %s", err))
		cancel()
		os.Exit(1)
	}

	authConfig := auth.Config{}
	if err := env.ParseWithOptions(&authConfig, env.Options{Prefix: envPrefixAuth}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s users configuration : %s", svcName, err))
		cancel()
		os.Exit(1)
	}

	uc, ucHandler, err := auth.Setup(authConfig)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer ucHandler.Close()

	logger.Info("Successfully connected to users grpc server " + ucHandler.Secure())

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

	svc := newService(collection, uc)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(otelgin.Middleware(svcName))
	router.Use(gin.Recovery())
	router.Use(helmet.Default())
	router.Use(ginprom.PromMiddleware(nil))
	router.Use(sloggin.New(logger))

	api.Endpoints(router, svc)

	httpServerConfig := server.Config{Port: defHTTPPort}
	if err := env.ParseWithOptions(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err.Error()))
		cancel()
		os.Exit(1)
	}
	hs := httpserver.New(ctx, cancel, svcName, httpServerConfig, router.Handler(), logger)

	g.Go(func() error {
		return hs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("%s service terminated: %s", svcName, err))
	}
}

func newService(collection *mongo.Collection, uc proto.UsersServiceClient) posts.Service {
	repo := repository.NewRepository(collection)

	return posts.NewService(repo, uc)
}

func connectDB(ctx context.Context, cfg config) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL).SetMonitor(otelmongo.NewMonitor()))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDB)

	if err := db.Client().Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := db.Collection(cfg.MongoColl)

	return collection, nil
}
