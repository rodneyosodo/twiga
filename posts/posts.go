// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package posts

import (
	"context"
	"encoding/json"
	"time"
)

type Comment struct {
	ID        string    `bson:"id,omitempty" json:"id"`
	Content   string    `bson:"content"      json:"content"`
	CreatedAt time.Time `bson:"created_at"   json:"created_at"`
	UpdateAt  time.Time `bson:"updated_at"   json:"updated_at"`
	UserID    string    `bson:"user_id"      json:"user_id"`
}

type CommentsPage struct {
	Page
	Comments []Comment `json:"comments"`
}

func (page CommentsPage) MarshalJSON() ([]byte, error) {
	type Alias CommentsPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Comments == nil {
		a.Comments = make([]Comment, 0)
	}

	return json.Marshal(a)
}

type Like struct {
	UserID    string    `bson:"user_id"    json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type LikesPage struct {
	Page
	Likes []Like `json:"likes"`
}

func (page LikesPage) MarshalJSON() ([]byte, error) {
	type Alias LikesPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Likes == nil {
		a.Likes = make([]Like, 0)
	}

	return json.Marshal(a)
}

type Share struct {
	ID        string    `bson:"id,omitempty" json:"id"`
	UserID    string    `bson:"user_id"      json:"user_id"`
	CreatedAt time.Time `bson:"created_at"   json:"created_at"`
}

type SharesPage struct {
	Page
	Shares []Share `json:"shares"`
}

func (page SharesPage) MarshalJSON() ([]byte, error) {
	type Alias SharesPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Shares == nil {
		a.Shares = make([]Share, 0)
	}

	return json.Marshal(a)
}

type Post struct {
	ID         string    `bson:"_id,omitempty"      json:"id"`
	Title      string    `bson:"title"              json:"title"`
	Content    string    `bson:"content"            json:"content"`
	Tags       []string  `bson:"tags"               json:"tags"`
	ImageURL   string    `bson:"image_url"          json:"image_url"`
	Visibility bool      `bson:"visibility"         json:"visibility"`
	CreatedAt  time.Time `bson:"created_at"         json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"         json:"updated_at"`
	UserID     string    `bson:"user_id"            json:"user_id"`
	Comments   []Comment `bson:"comments,omitempty" json:"comments,omitempty"`
	Likes      []Like    `bson:"likes,omitempty"    json:"likes,omitempty"`
	Shares     []Share   `bson:"shares,omitempty"   json:"shares,omitempty"`
}

func (post Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	a := struct {
		Alias
	}{
		Alias: Alias(post),
	}

	if a.Comments == nil {
		a.Comments = make([]Comment, 0)
	}
	if a.Likes == nil {
		a.Likes = make([]Like, 0)
	}
	if a.Shares == nil {
		a.Shares = make([]Share, 0)
	}

	return json.Marshal(a)
}

type Page struct {
	Total      uint64 `db:"total"                json:"total"`
	Offset     uint64 `db:"offset"               json:"offset"`
	Limit      uint64 `db:"limit"                json:"limit"`
	Tag        string `db:"tag,omitempty"        json:"tag,omitempty"`
	PostID     string `db:"post_id,omitempty"    json:"post_id,omitempty"`
	Visibility bool   `db:"visibility,omitempty" json:"visibility,omitempty"`
	UserID     string `db:"user_id,omitempty"    json:"user_id,omitempty"`
	CommentID  string `db:"comment_id,omitempty" json:"comment_id,omitempty"`
}

type PostsPage struct {
	Page
	Posts []Post `json:"posts"`
}

func (page PostsPage) MarshalJSON() ([]byte, error) {
	type Alias PostsPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Posts == nil {
		a.Posts = make([]Post, 0)
	}

	return json.Marshal(a)
}

//go:generate mockery --name Repository --output=./mocks --filename repository.go --quiet
type Repository interface { //nolint:interfacebloat
	Create(ctx context.Context, post Post) (Post, error)
	RetrieveByID(ctx context.Context, id string) (Post, error)
	RetrieveAll(ctx context.Context, page Page) (PostsPage, error)
	Update(ctx context.Context, post Post) (Post, error)
	UpdateContent(ctx context.Context, post Post) (Post, error)
	UpdateTags(ctx context.Context, post Post) (Post, error)
	UpdateImageURL(ctx context.Context, post Post) (Post, error)
	UpdateVisibility(ctx context.Context, post Post) (Post, error)
	Delete(ctx context.Context, id string) error

	CreateComment(ctx context.Context, postID string, comment Comment) (Comment, error)
	RetrieveCommentByID(ctx context.Context, id string) (Comment, error)
	RetrieveAllComments(ctx context.Context, page Page) (CommentsPage, error)
	UpdateComment(ctx context.Context, comment Comment) (Comment, error)
	DeleteComment(ctx context.Context, id string) error

	CreateLike(ctx context.Context, postID string, like Like) (Like, error)
	RetrieveAllLikes(ctx context.Context, page Page) (LikesPage, error)
	DeleteLike(ctx context.Context, postID, userID string) error

	CreateShare(ctx context.Context, postID string, share Share) (Share, error)
	RetrieveShareByID(ctx context.Context, id string) (Share, error)
	RetrieveAllShares(ctx context.Context, page Page) (SharesPage, error)
	DeleteShare(ctx context.Context, id string) error
}

//go:generate mockery --name Service --output=./mocks --filename service.go --quiet
type Service interface { //nolint:interfacebloat
	CreatePost(ctx context.Context, token string, post Post) (Post, error)
	RetrievePostByID(ctx context.Context, token string, id string) (Post, error)
	RetrieveAllPosts(ctx context.Context, token string, page Page) (PostsPage, error)
	UpdatePost(ctx context.Context, token string, post Post) (Post, error)
	UpdatePostContent(ctx context.Context, token string, post Post) (Post, error)
	UpdatePostTags(ctx context.Context, token string, post Post) (Post, error)
	UpdatePostImageURL(ctx context.Context, token string, post Post) (Post, error)
	UpdatePostVisibility(ctx context.Context, token string, post Post) (Post, error)
	DeletePost(ctx context.Context, token string, id string) error

	CreateComment(ctx context.Context, token string, postID string, comment Comment) (Comment, error)
	RetrieveCommentByID(ctx context.Context, token string, id string) (Comment, error)
	RetrieveAllComments(ctx context.Context, token string, page Page) (CommentsPage, error)
	UpdateComment(ctx context.Context, token string, comment Comment) (Comment, error)
	DeleteComment(ctx context.Context, roken, id string) error

	CreateLike(ctx context.Context, token string, postID string, like Like) (Like, error)
	RetrieveAllLikes(ctx context.Context, token string, page Page) (LikesPage, error)
	DeleteLike(ctx context.Context, token string, postID string) error

	CreateShare(ctx context.Context, token string, postID string, share Share) (Share, error)
	RetrieveAllShares(ctx context.Context, token string, page Page) (SharesPage, error)
	DeleteShare(ctx context.Context, token string, postID string) error
}
