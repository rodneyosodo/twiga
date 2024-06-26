// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package jaeger

import (
	"context"
	"errors"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	errNoURL                     = errors.New("URL is empty")
	errNoSvcName                 = errors.New("service Name is empty")
	errUnsupportedTraceURLScheme = errors.New("unsupported tracing url scheme")
)

func NewProvider(ctx context.Context, svcName string, jaegerUrl url.URL, fraction float64) (*trace.TracerProvider, error) {
	if jaegerUrl == (url.URL{}) {
		return nil, errNoURL
	}

	if svcName == "" {
		return nil, errNoSvcName
	}

	var client otlptrace.Client
	switch jaegerUrl.Scheme {
	case "http":
		client = otlptracehttp.NewClient(otlptracehttp.WithEndpoint(jaegerUrl.Host), otlptracehttp.WithURLPath(jaegerUrl.Path), otlptracehttp.WithInsecure())
	case "https":
		client = otlptracehttp.NewClient(otlptracehttp.WithEndpoint(jaegerUrl.Host), otlptracehttp.WithURLPath(jaegerUrl.Path))
	default:
		return nil, errUnsupportedTraceURLScheme
	}

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	attributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(svcName),
	}

	hostAttr, err := resource.New(ctx, resource.WithHost(), resource.WithOSDescription(), resource.WithContainer())
	if err != nil {
		return nil, err
	}
	attributes = append(attributes, hostAttr.Attributes()...)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(fraction)),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attributes...,
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}
