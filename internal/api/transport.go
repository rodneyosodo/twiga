// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractClientID(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return ""
	}

	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ")
	}

	return token
}
