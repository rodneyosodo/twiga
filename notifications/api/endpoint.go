// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"net/http"
	"strconv"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	iapi "github.com/rodneyosodo/twiga/internal/api"
	"github.com/rodneyosodo/twiga/notifications"
)

type entityfield string

const (
	contentField    entityfield = "content"
	tagsField       entityfield = "tags"
	imageField      entityfield = "image"
	visibilityField entityfield = "visibility"
	allFields       entityfield = "all"
)

func Endpoints(router *gin.Engine, svc notifications.Service) {
	router.POST("/notifications", createNotifications(svc))
	router.GET("/notifications", getNotifications(svc))
	router.GET("/notifications/:id", getNotification(svc))
	router.POST("/notifications/:id/read", readNotification(svc))
	router.POST("/notifications/read", readAllNotifications(svc))
	router.DELETE("/notifications/:id", deleteNotification(svc))

	router.POST("/settings", createSetting(svc))
	router.GET("/settings", getSettings(svc))
	router.GET("/settings/:id", getSetting(svc))
	router.PUT("/settings/:id", updateSetting(svc))
	router.DELETE("/settings/:id", deleteSetting(svc))

	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
}

func createNotifications(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		var notification notifications.Notification
		if err := ctx.BindJSON(&notification); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		notification, err := svc.CreateNotification(ctx, token, notification)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, notification)
	}
}

func getNotification(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})
		}

		notification, err := svc.RetrieveNotification(ctx, token, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, notification)
	}
}

func getNotifications(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		strOffset := ctx.DefaultQuery("offset", "0")
		offset, err := strconv.ParseUint(strOffset, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		strLimit := ctx.DefaultQuery("limit", "10")
		limit, err := strconv.ParseUint(strLimit, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		category := ctx.DefaultQuery("category", "")
		strIsRead := ctx.DefaultQuery("is_read", "true")
		isRead, err := strconv.ParseBool(strIsRead)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		page := notifications.Page{
			Offset:   offset,
			Limit:    limit,
			UserID:   ctx.DefaultQuery("user_id", ""),
			Category: notifications.ToCategory(category),
			IsRead:   &isRead,
		}

		notifications, err := svc.RetrieveAllNotifications(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, notifications)
	}
}

func readNotification(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})

			return
		}
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})

			return
		}

		if err := svc.ReadNotification(ctx, token, id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Post removed",
		})
	}
}

func readAllNotifications(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		strOffset := ctx.DefaultQuery("offset", "0")
		offset, err := strconv.ParseUint(strOffset, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		strLimit := ctx.DefaultQuery("limit", "10")
		limit, err := strconv.ParseUint(strLimit, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		category := ctx.DefaultQuery("category", "")
		strIsRead := ctx.DefaultQuery("is_read", "true")
		isRead, err := strconv.ParseBool(strIsRead)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		page := notifications.Page{
			Offset:   offset,
			Limit:    limit,
			UserID:   ctx.DefaultQuery("user_id", ""),
			Category: notifications.ToCategory(category),
			IsRead:   &isRead,
		}

		err = svc.ReadAllNotifications(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "All notifications read",
		})
	}
}

func deleteNotification(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})

			return
		}
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})

			return
		}

		if err := svc.DeleteNotification(ctx, token, id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Post removed",
		})
	}
}

func createSetting(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		postID := ctx.Param("id")
		if postID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
		}

		var setting notifications.Setting
		if err := ctx.BindJSON(&setting); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		setting, err := svc.CreateSetting(ctx, token, setting)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, setting)
	}
}

func getSetting(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})
		}

		setting, err := svc.RetrieveSetting(ctx, token, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, setting)
	}
}

func getSettings(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		postID := ctx.Param("id")
		if postID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
		}

		strOffset := ctx.DefaultQuery("offset", "0")
		offset, err := strconv.ParseUint(strOffset, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		strLimit := ctx.DefaultQuery("limit", "10")
		limit, err := strconv.ParseUint(strLimit, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		category := ctx.DefaultQuery("category", "")
		strIsRead := ctx.DefaultQuery("is_read", "true")
		isRead, err := strconv.ParseBool(strIsRead)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		page := notifications.Page{
			Offset:   offset,
			Limit:    limit,
			UserID:   ctx.DefaultQuery("user_id", ""),
			Category: notifications.ToCategory(category),
			IsRead:   &isRead,
		}

		settings, err := svc.RetrieveAllSettings(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, settings)
	}
}

func updateSetting(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		var setting notifications.Setting
		if err := ctx.BindJSON(&setting); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		setting.ID = ctx.Param("id")

		err := svc.UpdateSetting(ctx, token, setting)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Setting updated",
		})
	}
}

func deleteSetting(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})
		}

		if err := svc.DeleteSetting(ctx, token, id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Comment removed",
		})
	}
}
