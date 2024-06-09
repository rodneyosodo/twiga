// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type security int

const (
	withoutTLS security = iota
	withTLS
	withmTLS

	buffSize = 10 * 1024 * 1024
)

var (
	errGrpcConnect = errors.New("failed to connect to grpc server")
	errGrpcClose   = errors.New("failed to close grpc connection")
)

type Config struct {
	URL          string        `env:"URL"             envDefault:""`
	Timeout      time.Duration `env:"TIMEOUT"         envDefault:"1s"`
	ClientCert   string        `env:"CLIENT_CERT"     envDefault:""`
	ClientKey    string        `env:"CLIENT_KEY"      envDefault:""`
	ServerCAFile string        `env:"SERVER_CA_CERTS" envDefault:""`
}

type Handler interface {
	Close() error
	Secure() string
	Connection() *grpc.ClientConn
}

type client struct {
	*grpc.ClientConn
	cfg    Config
	secure security
}

var _ Handler = (*client)(nil)

func newHandler(cfg Config) (Handler, error) {
	conn, secure, err := connect(cfg)
	if err != nil {
		return nil, err
	}

	return &client{
		ClientConn: conn,
		cfg:        cfg,
		secure:     secure,
	}, nil
}

func (c *client) Close() error {
	if err := c.ClientConn.Close(); err != nil {
		return errors.Join(errGrpcClose, err)
	}

	return nil
}

func (c *client) Connection() *grpc.ClientConn {
	return c.ClientConn
}

func (c *client) Secure() string {
	switch c.secure {
	case withTLS:
		return "with TLS"
	case withmTLS:
		return "with mTLS"
	case withoutTLS:
		fallthrough
	default:
		return "without TLS"
	}
}

func connect(cfg Config) (*grpc.ClientConn, security, error) {
	opts := []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	secure := withoutTLS
	tc := insecure.NewCredentials()

	if cfg.ServerCAFile != "" { //nolint:nestif
		tlsConfig := &tls.Config{}

		rootCA, err := os.ReadFile(cfg.ServerCAFile)
		if err != nil {
			return nil, secure, fmt.Errorf("failed to load root ca file: %w", err)
		}
		if len(rootCA) > 0 {
			capool := x509.NewCertPool()
			if !capool.AppendCertsFromPEM(rootCA) {
				return nil, secure, fmt.Errorf("failed to append root ca to tls.Config")
			}
			tlsConfig.RootCAs = capool
			secure = withTLS
		}

		if cfg.ClientCert != "" || cfg.ClientKey != "" {
			certificate, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
			if err != nil {
				return nil, secure, fmt.Errorf("failed to client certificate and key %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{certificate}
			secure = withmTLS
		}

		tc = credentials.NewTLS(tlsConfig)
	}

	opts = append(
		opts, grpc.WithTransportCredentials(tc),
		grpc.WithReadBufferSize(buffSize),
		grpc.WithWriteBufferSize(buffSize),
	)

	conn, err := grpc.NewClient(cfg.URL, opts...)
	if err != nil {
		return nil, secure, errors.Join(errGrpcConnect, err)
	}

	return conn, secure, nil
}
