// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package repository_test

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"

	"github.com/0x6flab/namegenerator"
	"github.com/gofrs/uuid"
	"github.com/rodneyosodo/twiga/posts"
	"github.com/rodneyosodo/twiga/posts/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var namegen = namegenerator.NewGenerator()

func TestCreatePost(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: generatePost(),
			err:  nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			saved, err := repo.Create(context.Background(), tc.post)
			if tc.err != nil {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, saved.ID)
			assert.NotEmpty(t, saved.CreatedAt)
			assert.NotEmpty(t, saved.UpdatedAt)
		})
	}
}

func TestRetrieveByID(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc     string
		id       string
		response posts.Post
		err      error
	}{
		{
			desc:     "valid id",
			id:       saved.ID,
			response: saved,
			err:      nil,
		},
		{
			desc:     "invalid id",
			id:       "invalid_id",
			response: posts.Post{},
			err:      errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.RetrieveByID(context.Background(), tc.id)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response, post)
		})
	}
}

func TestRetrieveAll(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	num := uint64(10)
	saved := make([]posts.Post, num)
	for i := range num {
		post := generatePost()
		post.Visibility = true
		post, err := repo.Create(context.Background(), post)
		assert.NoError(t, err)
		saved[i] = post
	}

	cases := []struct {
		desc     string
		page     posts.Page
		response posts.PostsPage
		err      error
	}{
		{
			desc: "valid page",
			page: posts.Page{
				Offset: 0,
				Limit:  10,
			},
			response: posts.PostsPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  10,
				},
				Posts: saved,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: posts.Page{
				Offset: num + 1,
				Limit:  10,
			},
			response: posts.PostsPage{
				Page: posts.Page{
					Total:  num,
					Offset: num + 1,
					Limit:  10,
				},
				Posts: []posts.Post{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: posts.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			response: posts.PostsPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  num + 1,
				},
				Posts: saved,
			},
			err: nil,
		},
		{
			desc: "for a specific author",
			page: posts.Page{
				UserID: saved[0].UserID,
				Offset: 0,
				Limit:  10,
			},
			response: posts.PostsPage{
				Page: posts.Page{
					Total:  1,
					Offset: 0,
					Limit:  10,
				},
				Posts: []posts.Post{saved[0]},
			},
			err: nil,
		},
		{
			desc: "for tag",
			page: posts.Page{
				Tag:    saved[0].Tags[0],
				Offset: 0,
				Limit:  10,
			},
			response: posts.PostsPage{
				Page: posts.Page{
					Total:  1,
					Offset: 0,
					Limit:  10,
				},
				Posts: []posts.Post{saved[0]},
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			postsPage, err := repo.RetrieveAll(context.Background(), tc.page)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response.Total, postsPage.Total)
			assert.Equal(t, tc.response.Offset, postsPage.Offset)
			assert.Equal(t, tc.response.Limit, postsPage.Limit)
			assert.Equal(t, len(tc.response.Posts), len(postsPage.Posts))
			assert.ElementsMatch(t, tc.response.Posts, postsPage.Posts)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: posts.Post{
				ID:         saved.ID,
				Title:      namegen.Generate(),
				Content:    strings.Repeat("b", 100),
				Tags:       namegen.GenerateMultiple(3),
				ImageURL:   "https://www.example.com/image.jpg",
				Visibility: rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			desc: "invalid post",
			post: posts.Post{
				ID: "invalid_id",
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.Update(context.Background(), tc.post)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.post.Title, post.Title)
			assert.Equal(t, tc.post.Content, post.Content)
			assert.Equal(t, tc.post.Tags, post.Tags)
			assert.Equal(t, tc.post.ImageURL, post.ImageURL)
			assert.Equal(t, tc.post.Visibility, post.Visibility)
		})
	}
}

func TestUpdateContent(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: posts.Post{
				ID:      saved.ID,
				Content: strings.Repeat("b", 100),
			},
			err: nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			desc: "invalid post",
			post: posts.Post{
				ID: "invalid_id",
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.UpdateContent(context.Background(), tc.post)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.post.Title, post.Title)
			assert.Equal(t, tc.post.Content, post.Content)
			assert.Equal(t, tc.post.Tags, post.Tags)
			assert.Equal(t, tc.post.ImageURL, post.ImageURL)
			assert.Equal(t, tc.post.Visibility, post.Visibility)
		})
	}
}

func TestUpdateTags(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: posts.Post{
				ID:   saved.ID,
				Tags: namegen.GenerateMultiple(3),
			},
			err: nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			desc: "invalid post",
			post: posts.Post{
				ID: "invalid_id",
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.UpdateTags(context.Background(), tc.post)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.post.Title, post.Title)
			assert.Equal(t, tc.post.Content, post.Content)
			assert.Equal(t, tc.post.Tags, post.Tags)
			assert.Equal(t, tc.post.ImageURL, post.ImageURL)
			assert.Equal(t, tc.post.Visibility, post.Visibility)
		})
	}
}

func TestUpdateImageURL(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: posts.Post{
				ID:       saved.ID,
				ImageURL: "https://www.example.com/image.jpg",
			},
			err: nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			desc: "invalid post",
			post: posts.Post{
				ID: "invalid_id",
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.Update(context.Background(), tc.post)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.post.Title, post.Title)
			assert.Equal(t, tc.post.Content, post.Content)
			assert.Equal(t, tc.post.Tags, post.Tags)
			assert.Equal(t, tc.post.ImageURL, post.ImageURL)
			assert.Equal(t, tc.post.Visibility, post.Visibility)
		})
	}
}

func TestUpdateVisibility(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		post posts.Post
		err  error
	}{
		{
			desc: "valid post",
			post: posts.Post{
				ID:         saved.ID,
				Visibility: rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc: "empty post",
			post: posts.Post{},
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		{
			desc: "invalid post",
			post: posts.Post{
				ID: "invalid_id",
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			post, err := repo.Update(context.Background(), tc.post)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.post.Title, post.Title)
			assert.Equal(t, tc.post.Content, post.Content)
			assert.Equal(t, tc.post.Tags, post.Tags)
			assert.Equal(t, tc.post.ImageURL, post.ImageURL)
			assert.Equal(t, tc.post.Visibility, post.Visibility)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "valid id",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "invalid id",
			id:   "invalid_id",
			err:  errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.Delete(context.Background(), tc.id)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCreateComment(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc    string
		postID  string
		comment posts.Comment
		err     error
	}{
		{
			desc:   "valid comment",
			postID: saved.ID,
			comment: posts.Comment{
				Content: strings.Repeat("a", 100),
				UserID:  uuid.Must(uuid.NewV4()).String(),
			},
			err: nil,
		},
		{
			desc:    "empty comment",
			postID:  saved.ID,
			comment: posts.Comment{},
			err:     errors.New("user ID is required"),
		},
		{
			desc:   "invalid post id",
			postID: "invalid_id",
			comment: posts.Comment{
				Content: strings.Repeat("a", 100),
				UserID:  uuid.Must(uuid.NewV4()).String(),
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			comment, err := repo.CreateComment(context.Background(), tc.postID, tc.comment)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, comment.ID)
			assert.NotEmpty(t, comment.CreatedAt)
			assert.Equal(t, tc.comment.Content, comment.Content)

			post, err := repo.RetrieveByID(context.Background(), tc.postID)
			assert.NoError(t, err)

			assert.Len(t, post.Comments, 1)
		})
	}
}

func TestRetrieveComment(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	comment, err := repo.CreateComment(context.Background(), saved.ID, posts.Comment{
		Content: strings.Repeat("a", 100),
		UserID:  uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	_, err = repo.CreateComment(context.Background(), saved.ID, posts.Comment{
		Content: strings.Repeat("b", 100),
		UserID:  uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	cases := []struct {
		desc      string
		commentID string
		response  posts.Comment
		err       error
	}{
		{
			desc:      "valid comment",
			commentID: comment.ID,
			response:  comment,
			err:       nil,
		},
		{
			desc:      "invalid comment id",
			commentID: "invalid_id",
			response:  posts.Comment{},
			err:       errors.New("mongo: no documents in result"),
		},
		{
			desc:      "empty comment",
			commentID: "",
			response:  posts.Comment{},
			err:       errors.New("mongo: no documents in result"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			comment, err := repo.RetrieveCommentByID(context.Background(), tc.commentID)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response, comment)
		})
	}
}

func TestRetrieveAllComments(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	post := generatePost()
	post.Visibility = true
	post, err := repo.Create(context.Background(), post)
	assert.NoError(t, err)

	num := uint64(10)
	saved := make([]posts.Comment, num)
	for i := range num {
		comment, err := repo.CreateComment(context.Background(), post.ID, posts.Comment{
			Content: strings.Repeat(namegen.Generate(), 10),
			UserID:  uuid.Must(uuid.NewV4()).String(),
		})
		assert.NoError(t, err)

		saved[i] = comment
	}

	cases := []struct { //nolint:dupl
		desc     string
		page     posts.Page
		response posts.CommentsPage
		err      error
	}{
		{
			desc: "valid page",
			page: posts.Page{
				Offset: 0,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.CommentsPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  10,
				},
				Comments: saved,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: posts.Page{
				Offset: num + 1,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.CommentsPage{
				Page: posts.Page{
					Total:  num,
					Offset: num + 1,
					Limit:  10,
				},
				Comments: []posts.Comment{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: posts.Page{
				Offset: 0,
				Limit:  num + 1,
				PostID: post.ID,
			},
			response: posts.CommentsPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  num + 1,
				},
				Comments: saved,
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			commentsPage, err := repo.RetrieveAllComments(context.Background(), tc.page)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response.Total, commentsPage.Total)
			assert.Equal(t, tc.response.Offset, commentsPage.Offset)
			assert.Equal(t, tc.response.Limit, commentsPage.Limit)
			assert.Equal(t, len(tc.response.Comments), len(commentsPage.Comments))
			assert.ElementsMatch(t, tc.response.Comments, commentsPage.Comments)
		})
	}
}

func TestUpdateComment(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	comment, err := repo.CreateComment(context.Background(), saved.ID, posts.Comment{
		Content: strings.Repeat("a", 100),
		UserID:  uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	cases := []struct {
		desc    string
		comment posts.Comment
		err     error
	}{
		{
			desc: "valid comment",
			comment: posts.Comment{
				ID:      comment.ID,
				Content: strings.Repeat("b", 100),
			},
			err: nil,
		},
		{
			desc:    "empty comment",
			comment: posts.Comment{},
			err:     errors.New("comment not found"),
		},
		{
			desc: "invalid comment",
			comment: posts.Comment{
				ID: "invalid_id",
			},
			err: errors.New("comment not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			comment, err := repo.UpdateComment(context.Background(), tc.comment)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.comment.Content, comment.Content)

			post, err := repo.RetrieveByID(context.Background(), saved.ID)
			assert.NoError(t, err)

			assert.Len(t, post.Comments, 1)
			assert.Equal(t, tc.comment.Content, post.Comments[0].Content)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	comment, err := repo.CreateComment(context.Background(), saved.ID, posts.Comment{
		Content: strings.Repeat("a", 100),
		UserID:  uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	cases := []struct {
		desc    string
		comment posts.Comment
		err     error
	}{
		{
			desc: "valid comment",
			comment: posts.Comment{
				ID: comment.ID,
			},
			err: nil,
		},
		{
			desc:    "empty comment",
			comment: posts.Comment{},
			err:     errors.New("comment not found"),
		},
		{
			desc: "invalid comment",
			comment: posts.Comment{
				ID: "invalid_id",
			},
			err: errors.New("comment not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.DeleteComment(context.Background(), tc.comment.ID)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)

			post, err := repo.RetrieveByID(context.Background(), saved.ID)
			assert.NoError(t, err)

			assert.Len(t, post.Comments, 0)
		})
	}
}

func TestCreateLike(t *testing.T) { //nolint:dupl
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc   string
		postID string
		like   posts.Like
		err    error
	}{
		{
			desc:   "valid like",
			postID: saved.ID,
			like: posts.Like{
				UserID: uuid.Must(uuid.NewV4()).String(),
			},
			err: nil,
		},
		{
			desc:   "empty like",
			postID: saved.ID,
			like:   posts.Like{},
			err:    errors.New("user ID is required"),
		},
		{
			desc:   "invalid post id",
			postID: "invalid_id",
			like: posts.Like{
				UserID: uuid.Must(uuid.NewV4()).String(),
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			like, err := repo.CreateLike(context.Background(), tc.postID, tc.like)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, like.UserID)
			assert.NotEmpty(t, like.CreatedAt)

			post, err := repo.RetrieveByID(context.Background(), tc.postID)
			assert.NoError(t, err)

			assert.Len(t, post.Likes, 1)
		})
	}
}

func TestRetrieveAllLikes(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	post := generatePost()
	post.Visibility = true
	post, err := repo.Create(context.Background(), post)
	assert.NoError(t, err)

	num := uint64(10)
	saved := make([]posts.Like, num)
	for i := range num {
		like, err := repo.CreateLike(context.Background(), post.ID, posts.Like{
			UserID: uuid.Must(uuid.NewV4()).String(),
		})
		assert.NoError(t, err)

		saved[i] = like
	}

	cases := []struct { //nolint:dupl
		desc     string
		page     posts.Page
		response posts.LikesPage
		err      error
	}{
		{
			desc: "valid page",
			page: posts.Page{
				Offset: 0,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.LikesPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  10,
				},
				Likes: saved,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: posts.Page{
				Offset: num + 1,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.LikesPage{
				Page: posts.Page{
					Total:  num,
					Offset: num + 1,
					Limit:  10,
				},
				Likes: []posts.Like{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: posts.Page{
				Offset: 0,
				Limit:  num + 1,
				PostID: post.ID,
			},
			response: posts.LikesPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  num + 1,
				},
				Likes: saved,
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			likesPage, err := repo.RetrieveAllLikes(context.Background(), tc.page)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response.Total, likesPage.Total)
			assert.Equal(t, tc.response.Offset, likesPage.Offset)
			assert.Equal(t, tc.response.Limit, likesPage.Limit)
			assert.Equal(t, len(tc.response.Likes), len(likesPage.Likes))
			assert.ElementsMatch(t, tc.response.Likes, likesPage.Likes)
		})
	}
}

func TestDeleteLike(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	like, err := repo.CreateLike(context.Background(), saved.ID, posts.Like{
		UserID: uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	cases := []struct {
		desc string
		like posts.Like
		err  error
	}{
		{
			desc: "valid like",
			like: posts.Like{
				UserID: like.UserID,
			},
			err: nil,
		},
		{
			desc: "empty like",
			like: posts.Like{},
			err:  errors.New("like not found"),
		},
		{
			desc: "invalid like",
			like: posts.Like{
				UserID: "invalid_id",
			},
			err: errors.New("like not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.DeleteLike(context.Background(), tc.like.UserID)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)

			post, err := repo.RetrieveByID(context.Background(), saved.ID)
			assert.NoError(t, err)

			assert.Len(t, post.Likes, 0)
		})
	}
}

func TestCreateShare(t *testing.T) { //nolint:dupl
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	cases := []struct {
		desc   string
		postID string
		share  posts.Share
		err    error
	}{
		{
			desc:   "valid share",
			postID: saved.ID,
			share: posts.Share{
				UserID: uuid.Must(uuid.NewV4()).String(),
			},
			err: nil,
		},
		{
			desc:   "empty share",
			postID: saved.ID,
			share:  posts.Share{},
			err:    errors.New("user ID is required"),
		},
		{
			desc:   "invalid post id",
			postID: "invalid_id",
			share: posts.Share{
				UserID: uuid.Must(uuid.NewV4()).String(),
			},
			err: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			share, err := repo.CreateShare(context.Background(), tc.postID, tc.share)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, share.UserID)
			assert.NotEmpty(t, share.CreatedAt)

			post, err := repo.RetrieveByID(context.Background(), tc.postID)
			assert.NoError(t, err)

			assert.Len(t, post.Shares, 1)
		})
	}
}

func TestRetrieveAllShares(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	post := generatePost()
	post.Visibility = true
	post, err := repo.Create(context.Background(), post)
	assert.NoError(t, err)

	num := uint64(10)
	saved := make([]posts.Share, num)
	for i := range num {
		share, err := repo.CreateShare(context.Background(), post.ID, posts.Share{
			UserID: uuid.Must(uuid.NewV4()).String(),
		})
		assert.NoError(t, err)

		saved[i] = share
	}

	cases := []struct { //nolint:dupl
		desc     string
		page     posts.Page
		response posts.SharesPage
		err      error
	}{
		{
			desc: "valid page",
			page: posts.Page{
				Offset: 0,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.SharesPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  10,
				},
				Shares: saved,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: posts.Page{
				Offset: num + 1,
				Limit:  10,
				PostID: post.ID,
			},
			response: posts.SharesPage{
				Page: posts.Page{
					Total:  num,
					Offset: num + 1,
					Limit:  10,
				},
				Shares: []posts.Share{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: posts.Page{
				Offset: 0,
				Limit:  num + 1,
				PostID: post.ID,
			},
			response: posts.SharesPage{
				Page: posts.Page{
					Total:  num,
					Offset: 0,
					Limit:  num + 1,
				},
				Shares: saved,
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			sharesPage, err := repo.RetrieveAllShares(context.Background(), tc.page)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.response.Total, sharesPage.Total)
			assert.Equal(t, tc.response.Offset, sharesPage.Offset)
			assert.Equal(t, tc.response.Limit, sharesPage.Limit)
			assert.Equal(t, len(tc.response.Shares), len(sharesPage.Shares))
			assert.ElementsMatch(t, tc.response.Shares, sharesPage.Shares)
		})
	}
}

func TestDeleteShare(t *testing.T) {
	t.Cleanup(func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		assert.NoError(t, err)
	})

	repo := repository.NewRepository(collection)

	saved, err := repo.Create(context.Background(), generatePost())
	assert.NoError(t, err)

	share, err := repo.CreateShare(context.Background(), saved.ID, posts.Share{
		UserID: uuid.Must(uuid.NewV4()).String(),
	})
	assert.NoError(t, err)

	cases := []struct {
		desc  string
		share posts.Share
		err   error
	}{
		{
			desc: "valid share",
			share: posts.Share{
				UserID: share.UserID,
			},
			err: nil,
		},
		{
			desc:  "empty share",
			share: posts.Share{},
			err:   errors.New("share not found"),
		},
		{
			desc: "invalid share",
			share: posts.Share{
				UserID: "invalid_id",
			},
			err: errors.New("share not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.DeleteShare(context.Background(), tc.share.UserID)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)

			post, err := repo.RetrieveByID(context.Background(), saved.ID)
			assert.NoError(t, err)

			assert.Len(t, post.Shares, 0)
		})
	}
}

func generatePost() posts.Post {
	return posts.Post{
		Title:      namegen.Generate(),
		Content:    strings.Repeat("a", 100),
		Tags:       namegen.GenerateMultiple(3),
		ImageURL:   "https://www.example.com/image.jpg",
		Visibility: rand.Intn(2) == 0,
		UserID:     uuid.Must(uuid.NewV4()).String(),
	}
}
