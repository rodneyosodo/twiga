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
	"github.com/chenjiandongx/ginprom"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	"github.com/grafana/loki-client-go/loki"
	"github.com/redis/go-redis/v9"
	"github.com/rodneyosodo/twiga/internal/auth"
	"github.com/rodneyosodo/twiga/internal/cache"
	"github.com/rodneyosodo/twiga/internal/events"
	"github.com/rodneyosodo/twiga/internal/events/rabbitmq"
	"github.com/rodneyosodo/twiga/internal/jaeger"
	"github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/internal/server"
	httpserver "github.com/rodneyosodo/twiga/internal/server/http"
	"github.com/rodneyosodo/twiga/notifications"
	"github.com/rodneyosodo/twiga/notifications/api"
	"github.com/rodneyosodo/twiga/notifications/consumer"
	"github.com/rodneyosodo/twiga/notifications/repository"
	sloggin "github.com/samber/slog-gin"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"golang.org/x/sync/errgroup"
)

const (
	svcName       = "notifications"
	defHTTPPort   = "6002"
	envPrefixHTTP = "TWIGA_NOTIFICATIONS_HTTP_"
	envPrefixAuth = "TWIGA_USERS_GRPC_"
	envPrefixDB   = "TWIGA_NOTIFICATIONS_DB_"
	defDB         = "notifications"
)

type config struct {
	LogLevel         string        `env:"TWIGA_NOTIFICATIONS_LOG_LEVEL" envDefault:"info"`
	JaegerURL        url.URL       `env:"TWIGA_JAEGER_URL"              envDefault:"http://localhost:14268"`
	TraceRatio       float64       `env:"TWIGA_JAEGER_TRACE_RATIO"      envDefault:"1.0"`
	ESURL            string        `env:"TWIGA_ES_URL"                  envDefault:"amqp://twiga:twiga@localhost:5672/"`
	CacheURL         string        `env:"TWIGA_CACHE_URL"               envDefault:"redis://localhost:6379/0"`
	CacheKeyDuration time.Duration `env:"TWIGA_CACHE_KEY_DURATION"      envDefault:"10m"`
	LokiURL          string        `env:"TWIGA_LOKI_URL"                envDefault:"http://localhost:3100/loki/api/v1/push"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err.Error())
	}

	config, err := loki.NewDefaultConfig(cfg.LokiURL)
	if err != nil {
		log.Fatalf("failed to create loki config: %s", err.Error())
	}
	config.TenantID = svcName
	config.EncodeJson = true
	client, err := loki.New(config)
	if err != nil {
		log.Fatalf("failed to create loki client: %s", err.Error())
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		log.Fatalf("failed to parse log level: %s", err.Error())
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	lokiHandler := slogloki.Option{Level: level, Client: client}.NewLokiHandler()

	logger := slog.New(slogmulti.Fanout(logHandler, lokiHandler))

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

	dbConfig := postgres.Config{Name: defDB}
	if err := env.ParseWithOptions(&dbConfig, env.Options{Prefix: envPrefixDB}); err != nil {
		logger.Error(err.Error())
	}

	db, err := postgres.Setup(dbConfig, *repository.Migration())
	if err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}
	defer db.Close()

	authConfig := auth.Config{}
	if err := env.ParseWithOptions(&authConfig, env.Options{Prefix: envPrefixAuth}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s users configuration : %s", svcName, err))
		cancel()
		os.Exit(1)
	}

	uc, ucHandler, err := auth.Setup(authConfig)
	if err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}
	defer ucHandler.Close()

	logger.Info("Successfully connected to users grpc server " + ucHandler.Secure())

	repo := repository.NewRepository(db)
	cacher, err := connectToCache(ctx, cfg.CacheURL, cfg.CacheKeyDuration)
	if err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	svc := notifications.NewService(repo, uc, cacher)

	pubsub, err := rabbitmq.NewPubSub(cfg.ESURL, logger)
	if err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}
	subConfig := events.SubscriberConfig{
		ID:      svcName,
		Topic:   rabbitmq.SubjectAllEvents,
		Handler: consumer.NewEventHandler(svc),
	}
	if err := pubsub.Subscribe(ctx, subConfig); err != nil {
		logger.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(otelgin.Middleware(svcName, otelgin.WithTracerProvider(tp)))
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

func connectToCache(ctx context.Context, url string, duration time.Duration) (cache.Cacher, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return cache.NewCache(redis.NewClient(opts), duration), nil
}
