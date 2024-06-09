// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package producer

import (
	"context"
	"encoding/json"

	"github.com/rodneyosodo/twiga/internal/events"
	"github.com/rodneyosodo/twiga/posts"
)

var _ posts.Service = (*eventStore)(nil)

type eventStore struct {
	publisher events.Publisher
	svc       posts.Service
}

func NewEventStore(pub events.Publisher, svc posts.Service) posts.Service {
	return &eventStore{
		publisher: pub,
		svc:       svc,
	}
}

func (e *eventStore) CreatePost(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.CreatePost(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.created", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) RetrievePostByID(ctx context.Context, token string, id string) (posts.Post, error) {
	return e.svc.RetrievePostByID(ctx, token, id)
}

func (e *eventStore) RetrieveAllPosts(ctx context.Context, token string, page posts.Page) (posts.PostsPage, error) {
	return e.svc.RetrieveAllPosts(ctx, token, page)
}

func (e *eventStore) UpdatePost(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.UpdatePost(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.updated", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) UpdatePostContent(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.UpdatePostContent(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.updated.content", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) UpdatePostTags(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.UpdatePostTags(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.updated.tags", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) UpdatePostImageURL(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.UpdatePostImageURL(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.updated.image", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) UpdatePostVisibility(ctx context.Context, token string, post posts.Post) (posts.Post, error) {
	post, err := e.svc.UpdatePostVisibility(ctx, token, post)
	if err != nil {
		return posts.Post{}, err
	}

	if err := Publish(ctx, e.publisher, "posts.updated.visibility", post); err != nil {
		return post, err
	}

	return post, nil
}

func (e *eventStore) DeletePost(ctx context.Context, token string, id string) error {
	if err := e.svc.DeletePost(ctx, token, id); err != nil {
		return err
	}

	if err := Publish(ctx, e.publisher, "posts.deleted", posts.Post{ID: id}); err != nil {
		return err
	}

	return nil
}

func (e *eventStore) CreateComment(ctx context.Context, token string, postID string, comment posts.Comment) (posts.Comment, error) {
	comment, err := e.svc.CreateComment(ctx, token, postID, comment)
	if err != nil {
		return posts.Comment{}, err
	}

	if err := Publish(ctx, e.publisher, "comments.created", comment); err != nil {
		return comment, err
	}

	return comment, nil
}

func (e *eventStore) RetrieveCommentByID(ctx context.Context, token string, id string) (posts.Comment, error) {
	return e.svc.RetrieveCommentByID(ctx, token, id)
}

func (e *eventStore) RetrieveAllComments(ctx context.Context, token string, page posts.Page) (posts.CommentsPage, error) {
	return e.svc.RetrieveAllComments(ctx, token, page)
}

func (e *eventStore) UpdateComment(ctx context.Context, token string, comment posts.Comment) (posts.Comment, error) {
	comment, err := e.svc.UpdateComment(ctx, token, comment)
	if err != nil {
		return posts.Comment{}, err
	}

	if err := Publish(ctx, e.publisher, "comments.updated", comment); err != nil {
		return comment, err
	}

	return comment, nil
}

func (e *eventStore) DeleteComment(ctx context.Context, token string, id string) error {
	if err := e.svc.DeleteComment(ctx, token, id); err != nil {
		return err
	}

	if err := Publish(ctx, e.publisher, "comments.deleted", posts.Comment{ID: id}); err != nil {
		return err
	}

	return nil
}

func (e *eventStore) CreateLike(ctx context.Context, token string, postID string, like posts.Like) (posts.Like, error) {
	like, err := e.svc.CreateLike(ctx, token, postID, like)
	if err != nil {
		return posts.Like{}, err
	}

	if err := Publish(ctx, e.publisher, "likes.created", like); err != nil {
		return like, err
	}

	return like, nil
}

func (e *eventStore) RetrieveAllLikes(ctx context.Context, token string, page posts.Page) (posts.LikesPage, error) {
	return e.svc.RetrieveAllLikes(ctx, token, page)
}

func (e *eventStore) DeleteLike(ctx context.Context, token, postID string) error {
	if err := e.svc.DeleteLike(ctx, token, postID); err != nil {
		return err
	}

	if err := Publish(ctx, e.publisher, "likes.deleted", posts.Post{ID: postID}); err != nil {
		return err
	}

	return nil
}

func (e *eventStore) CreateShare(ctx context.Context, token string, postID string, share posts.Share) (posts.Share, error) {
	share, err := e.svc.CreateShare(ctx, token, postID, share)
	if err != nil {
		return posts.Share{}, err
	}

	if err := Publish(ctx, e.publisher, "shares.created", share); err != nil {
		return share, err
	}

	return share, nil
}

func (e *eventStore) RetrieveAllShares(ctx context.Context, token string, page posts.Page) (posts.SharesPage, error) {
	return e.svc.RetrieveAllShares(ctx, token, page)
}

func (e *eventStore) DeleteShare(ctx context.Context, token string, id string) error {
	if err := e.svc.DeleteShare(ctx, token, id); err != nil {
		return err
	}

	if err := Publish(ctx, e.publisher, "shares.deleted", posts.Share{ID: id}); err != nil {
		return err
	}

	return nil
}

type entity interface {
	posts.Post | posts.Comment | posts.Like | posts.Share
}

func Publish[E entity](ctx context.Context, pub events.Publisher, topic string, entity E) error {
	jsonEntity, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(jsonEntity, &payload); err != nil {
		return err
	}

	return pub.Publish(ctx, topic, payload)
}
