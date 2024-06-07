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
	"github.com/rodneyosodo/twiga/posts"
)

type entityfield string

const (
	contentField    entityfield = "content"
	tagsField       entityfield = "tags"
	imageField      entityfield = "image"
	visibilityField entityfield = "visibility"
	allFields       entityfield = "all"
)

func Endpoints(router *gin.Engine, svc posts.Service) {
	router.POST("/posts", createPost(svc))
	router.GET("/posts", getPosts(svc))
	router.GET("/posts/:id", getPost(svc))
	router.PUT("/posts/:id", updatePost(svc, allFields))
	router.PATCH("/posts/:id/content", updatePost(svc, contentField))
	router.PATCH("/posts/:id/tags", updatePost(svc, tagsField))
	router.PATCH("/posts/:id/image", updatePost(svc, imageField))
	router.PATCH("/posts/:id/visibility", updatePost(svc, visibilityField))
	router.DELETE("/posts/:id", deletePost(svc))

	router.POST("/posts/:id/comments", createComment(svc))
	router.GET("/posts/:id/comments", getComments(svc))
	router.GET("/posts/comments/:id", getComment(svc))
	router.PUT("/posts/comments/:id", updateComment(svc))
	router.DELETE("/posts/comments/:id", deleteComment(svc))

	router.POST("/posts/:id/like", likePost(svc))
	router.GET("/posts/:id/likes", getLikes(svc))
	router.DELETE("/posts/:id/unlike", deleteLike(svc))

	router.POST("/posts/:id/share", sharePost(svc))
	router.GET("/posts/:id/shares", getShares(svc))
	router.DELETE("/posts/:id/unshare", deleteShare(svc))

	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
}

func createPost(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		var post posts.Post
		if err := ctx.BindJSON(&post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		post, err := svc.CreatePost(ctx, token, post)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, post)
	}
}

func getPost(svc posts.Service) func(ctx *gin.Context) {
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

		post, err := svc.RetrievePostByID(ctx, token, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, post)
	}
}

func getPosts(svc posts.Service) func(ctx *gin.Context) {
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
		strVisibility := ctx.DefaultQuery("visibility", "true")
		visibility, err := strconv.ParseBool(strVisibility)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		page := posts.Page{
			Offset:     offset,
			Limit:      limit,
			Visibility: visibility,
			UserID:     ctx.DefaultQuery("user_id", ""),
		}

		posts, err := svc.RetrieveAllPosts(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, posts)
	}
}

func updatePost(svc posts.Service, field entityfield) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		var post posts.Post
		if err := ctx.BindJSON(&post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		post.ID = ctx.Param("id")

		switch field {
		case contentField:
			if post.Content == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Content is required",
				})
			}
		case tagsField:
			if len(post.Tags) == 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Tags are required",
				})
			}
		case imageField:
			if post.ImageURL == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Image is required",
				})
			}
		case visibilityField:
			if post.Visibility == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Visibility is required",
				})
			}
		}

		var updated posts.Post
		var err error
		switch field {
		case contentField:
			updated, err = svc.UpdatePostContent(ctx, token, post)
		case tagsField:
			updated, err = svc.UpdatePostTags(ctx, token, post)
		case imageField:
			updated, err = svc.UpdatePostImageURL(ctx, token, post)
		case visibilityField:
			updated, err = svc.UpdatePostVisibility(ctx, token, post)
		case allFields:
			updated, err = svc.UpdatePost(ctx, token, post)
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, updated)
	}
}

func deletePost(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		if err := svc.DeletePost(ctx, token, ctx.Param("id")); err != nil {
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

func createComment(svc posts.Service) func(ctx *gin.Context) { //nolint:dupl
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

		var comment posts.Comment
		if err := ctx.BindJSON(&comment); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		comment, err := svc.CreateComment(ctx, token, postID, comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, comment)
	}
}

func getComment(svc posts.Service) func(ctx *gin.Context) {
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

		comment, err := svc.RetrieveCommentByID(ctx, token, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, comment)
	}
}

func getComments(svc posts.Service) func(ctx *gin.Context) {
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

		page := posts.Page{
			Offset: offset,
			Limit:  limit,
			UserID: ctx.DefaultQuery("user_id", ""),
			PostID: postID,
		}

		comments, err := svc.RetrieveAllComments(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, comments)
	}
}

func updateComment(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		var comment posts.Comment
		if err := ctx.BindJSON(&comment); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		comment.ID = ctx.Param("id")

		updated, err := svc.UpdateComment(ctx, token, comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, updated)
	}
}

func deleteComment(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		if err := svc.DeleteComment(ctx, token, ctx.Param("id")); err != nil {
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

func likePost(svc posts.Service) func(ctx *gin.Context) { //nolint:dupl
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

		var like posts.Like
		if err := ctx.BindJSON(&like); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		like, err := svc.CreateLike(ctx, token, postID, like)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, like)
	}
}

func getLikes(svc posts.Service) func(ctx *gin.Context) {
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

		page := posts.Page{
			Offset: offset,
			Limit:  limit,
			UserID: ctx.DefaultQuery("user_id", ""),
		}

		likes, err := svc.RetrieveAllLikes(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, likes)
	}
}

func deleteLike(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		if err := svc.DeleteLike(ctx, token, ctx.Param("id")); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Like removed",
		})
	}
}

func sharePost(svc posts.Service) func(ctx *gin.Context) { //nolint:dupl
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

		var share posts.Share
		if err := ctx.BindJSON(&share); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		share, err := svc.CreateShare(ctx, token, postID, share)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusCreated, share)
	}
}

func getShares(svc posts.Service) func(ctx *gin.Context) {
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

		page := posts.Page{
			Offset: offset,
			Limit:  limit,
			UserID: ctx.DefaultQuery("user_id", ""),
		}

		shares, err := svc.RetrieveAllShares(ctx, token, page)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, shares)
	}
}

func deleteShare(svc posts.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := iapi.GinExtractToken(ctx)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
		}

		if err := svc.DeleteShare(ctx, token, ctx.Param("id")); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "Share removed",
		})
	}
}
