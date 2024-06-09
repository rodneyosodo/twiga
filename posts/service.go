// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package posts

import (
	"context"
	"errors"

	"github.com/rodneyosodo/twiga/internal/cache"
	"github.com/rodneyosodo/twiga/users/proto"
)

var _ Service = (*service)(nil)

type service struct {
	repo   Repository
	users  proto.UsersServiceClient
	cacher cache.Cacher
}

func NewService(repo Repository, users proto.UsersServiceClient, cacher cache.Cacher) Service {
	return &service{
		repo:   repo,
		users:  users,
		cacher: cacher,
	}
}

func (s *service) CreatePost(ctx context.Context, token string, post Post) (Post, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Post{}, err
	}
	post.UserID = userID

	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.Create(ctx, post)
}

func (s *service) RetrievePostByID(ctx context.Context, token string, id string) (saved Post, err error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return Post{}, err
	}
	if cached := s.cacher.Get(ctx, id); cached != nil {
		saved, ok := cached.(Post)
		if ok {
			return saved, nil
		}
	}

	return s.repo.RetrieveByID(ctx, id)
}

func (s *service) RetrieveAllPosts(ctx context.Context, token string, page Page) (PostsPage, error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return PostsPage{}, err
	}

	return s.repo.RetrieveAll(ctx, page)
}

func (s *service) UpdatePost(ctx context.Context, token string, post Post) (updated Post, err error) {
	if err := s.authorize(ctx, token, post.ID); err != nil {
		return Post{}, err
	}

	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.Update(ctx, post)
}

func (s *service) authorize(ctx context.Context, token, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}

	if cached := s.cacher.Get(ctx, id); cached != nil {
		saved, ok := cached.(Post)
		if ok {
			if saved.ID == userID {
				return nil
			}
		}
	}

	saved, err := s.repo.RetrieveByID(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return nil
}

func (s *service) UpdatePostContent(ctx context.Context, token string, post Post) (updated Post, err error) {
	if err := s.authorize(ctx, token, post.ID); err != nil {
		return Post{}, err
	}
	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.UpdateContent(ctx, post)
}

func (s *service) UpdatePostTags(ctx context.Context, token string, post Post) (updated Post, err error) {
	if err := s.authorize(ctx, token, post.ID); err != nil {
		return Post{}, err
	}
	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.UpdateTags(ctx, post)
}

func (s *service) UpdatePostImageURL(ctx context.Context, token string, post Post) (updated Post, err error) {
	if err := s.authorize(ctx, token, post.ID); err != nil {
		return Post{}, err
	}
	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.UpdateImageURL(ctx, post)
}

func (s *service) UpdatePostVisibility(ctx context.Context, token string, post Post) (updated Post, err error) {
	if err := s.authorize(ctx, token, post.ID); err != nil {
		return Post{}, err
	}
	defer func() {
		if err == nil {
			if err = s.cacher.Add(ctx, post.ID, post); err != nil {
				err = errors.New("failed to cache post")
			}
		}
	}()

	return s.repo.UpdateVisibility(ctx, post)
}

func (s *service) DeletePost(ctx context.Context, token string, id string) error {
	if err := s.authorize(ctx, token, id); err != nil {
		return err
	}
	if err := s.cacher.Remove(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) CreateComment(ctx context.Context, token string, postID string, comment Comment) (Comment, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Comment{}, err
	}
	comment.UserID = userID

	return s.repo.CreateComment(ctx, postID, comment)
}

func (s *service) RetrieveCommentByID(ctx context.Context, token string, id string) (Comment, error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return Comment{}, err
	}

	return s.repo.RetrieveCommentByID(ctx, id)
}

func (s *service) RetrieveAllComments(ctx context.Context, token string, page Page) (CommentsPage, error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return CommentsPage{}, err
	}

	return s.repo.RetrieveAllComments(ctx, page)
}

func (s *service) UpdateComment(ctx context.Context, token string, comment Comment) (Comment, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Comment{}, err
	}
	saved, err := s.repo.RetrieveCommentByID(ctx, comment.ID)
	if err != nil {
		return Comment{}, err
	}
	if saved.UserID != userID {
		return Comment{}, errors.New("unauthorized")
	}

	return s.repo.UpdateComment(ctx, comment)
}

func (s *service) DeleteComment(ctx context.Context, token string, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveCommentByID(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteComment(ctx, id)
}

func (s *service) CreateLike(ctx context.Context, token string, postID string, like Like) (Like, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Like{}, err
	}
	like.UserID = userID

	return s.repo.CreateLike(ctx, postID, like)
}

func (s *service) RetrieveAllLikes(ctx context.Context, token string, page Page) (LikesPage, error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return LikesPage{}, err
	}

	return s.repo.RetrieveAllLikes(ctx, page)
}

func (s *service) DeleteLike(ctx context.Context, token, postID string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}

	return s.repo.DeleteLike(ctx, postID, userID)
}

func (s *service) CreateShare(ctx context.Context, token string, postID string, share Share) (Share, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Share{}, err
	}
	share.UserID = userID

	return s.repo.CreateShare(ctx, postID, share)
}

func (s *service) RetrieveAllShares(ctx context.Context, token string, page Page) (SharesPage, error) {
	if _, err := s.IdentifyUser(ctx, token); err != nil {
		return SharesPage{}, err
	}

	return s.repo.RetrieveAllShares(ctx, page)
}

func (s *service) DeleteShare(ctx context.Context, token string, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveShareByID(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteShare(ctx, id)
}

func (s *service) IdentifyUser(ctx context.Context, token string) (string, error) {
	resp, err := s.users.IdentifyUser(ctx, &proto.IdentifyUserRequest{Token: token})
	if err != nil {
		return "", err
	}

	return resp.GetId(), nil
}
