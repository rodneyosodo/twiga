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
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	bufferSize   = 1024
	defaultLimit = 100
	defaultSleep = 5 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Endpoints(router *gin.Engine, svc notifications.Service) {
	router.GET("/notifications", getNotifications(svc))
	router.GET("/notifications/:id", getNotification(svc))
	router.POST("/notifications/:id/read", readNotification(svc))
	router.POST("/notifications/read", readAllNotifications(svc))
	router.DELETE("/notifications/:id", deleteNotification(svc))

	router.GET("/ws", wsHandler(svc))

	router.GET("/version", iapi.GinVersion("notifications"))
	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
}

func wsHandler(svc notifications.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})

			return
		}

		category := ctx.DefaultQuery("category", "")
		strIsRead := ctx.DefaultQuery("is_read", "false")
		isRead, err := strconv.ParseBool(strIsRead)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)

			return
		}
		defer conn.Close()

		pm := notifications.Page{
			Offset:   0,
			Limit:    defaultLimit,
			Category: notifications.ToCategory(category),
			IsRead:   &isRead,
		}

		notifs, err := svc.RetrieveAllNotifications(ctx, token, pm)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		for notifs.Total > uint64(len(notifs.Notifications)) {
			pm.Offset += defaultLimit
			newNotifs, err := svc.RetrieveAllNotifications(ctx, token, pm)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})

				return
			}
			notifs.Notifications = append(notifs.Notifications, newNotifs.Notifications...)
		}

		for _, n := range notifs.Notifications {
			if err := conn.WriteJSON(n); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)

				return
			}
		}

		userID, err := svc.IdentifyUser(ctx, token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		for {
			newNotif := svc.GetNewNotification(ctx, userID)
			if newNotif.ID == "" {
				goto sleep
			}
			if err := conn.WriteJSON(newNotif); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)

				break
			}
		sleep:
			time.Sleep(defaultSleep)
		}
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
		strIsRead := ctx.DefaultQuery("is_read", "false")
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
