// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"encoding/json"
	"net/http"

	"go.szostok.io/version"
)

func Version(service string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/health+json")
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		info := version.Get()

		res := map[string]string{
			"service":     service,
			"status":      "pass",
			"version":     info.Version,
			"commit":      info.GitCommit,
			"commit_date": info.CommitDate,
			"build_date":  info.BuildDate,
			"go_version":  info.GoVersion,
			"compiler":    info.Compiler,
			"platform":    info.Platform,
		}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
