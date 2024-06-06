// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rodneyosodo/twiga/posts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ posts.Repository = (*repository)(nil)

type repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Collection) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, post posts.Post) (posts.Post, error) {
	post.ID = ""
	post.CreatedAt = time.Now().UTC().Round(time.Millisecond)
	post.UpdatedAt = time.Now().UTC().Round(time.Millisecond)

	saved, err := r.db.InsertOne(ctx, post)
	if err != nil {
		return posts.Post{}, err
	}
	id, ok := saved.InsertedID.(primitive.ObjectID)
	if !ok {
		return post, errors.New("could not get post ID")
	}
	post.ID = id.Hex()

	return post, nil
}

func (r *repository) RetrieveByID(ctx context.Context, id string) (posts.Post, error) {
	var post posts.Post
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, err
	}

	if err = r.db.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&post); err != nil {
		return post, err
	}

	return post, nil
}

func (r *repository) RetrieveAll(ctx context.Context, page posts.Page) (posts.PostsPage, error) {
	filter, opts := getFilter(page)

	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return posts.PostsPage{}, err
	}
	defer cursor.Close(ctx)

	var postsPage posts.PostsPage
	for cursor.Next(ctx) {
		var post posts.Post
		if err := cursor.Decode(&post); err != nil {
			return postsPage, err
		}
		postsPage.Posts = append(postsPage.Posts, post)
	}

	total, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return posts.PostsPage{}, err
	}
	postsPage.Total = uint64(total)
	postsPage.Offset = page.Offset
	postsPage.Limit = page.Limit

	return postsPage, nil
}

func (r *repository) Update(ctx context.Context, post posts.Post) (posts.Post, error) {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return post, err
	}
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: post.Title},
		{Key: "content", Value: post.Content},
		{Key: "tags", Value: post.Tags},
		{Key: "image_url", Value: post.ImageURL},
		{Key: "visibility", Value: post.Visibility},
		{Key: "updated_at", Value: time.Now().UTC().Round(time.Millisecond)},
	}}})
	if err != nil {
		return post, err
	}
	if result.MatchedCount == 0 {
		return post, errors.New("post not found")
	}

	return post, nil
}

func (r *repository) UpdateContent(ctx context.Context, post posts.Post) (posts.Post, error) {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return post, err
	}
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "content", Value: post.Content},
		{Key: "updated_at", Value: time.Now().UTC().Round(time.Millisecond)},
	}}})
	if err != nil {
		return post, err
	}
	if result.MatchedCount == 0 {
		return post, errors.New("post not found")
	}

	return post, nil
}

func (r *repository) UpdateTags(ctx context.Context, post posts.Post) (posts.Post, error) {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return post, err
	}
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "tags", Value: post.Tags},
		{Key: "updated_at", Value: time.Now().UTC().Round(time.Millisecond)},
	}}})
	if err != nil {
		return post, err
	}
	if result.MatchedCount == 0 {
		return post, errors.New("post not found")
	}

	return post, nil
}

func (r *repository) UpdateImageURL(ctx context.Context, post posts.Post) (posts.Post, error) {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return post, err
	}
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "image_url", Value: post.ImageURL},
		{Key: "updated_at", Value: time.Now().UTC().Round(time.Millisecond)},
	}}})
	if err != nil {
		return post, err
	}
	if result.MatchedCount == 0 {
		return post, errors.New("post not found")
	}

	return post, nil
}

func (r *repository) UpdateVisibility(ctx context.Context, post posts.Post) (posts.Post, error) {
	objID, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		return post, err
	}
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "visibility", Value: post.Visibility},
		{Key: "updated_at", Value: time.Now().UTC().Round(time.Millisecond)},
	}}})
	if err != nil {
		return post, err
	}
	if result.MatchedCount == 0 {
		return post, errors.New("post not found")
	}

	return post, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := r.db.DeleteOne(ctx, bson.D{{Key: "_id", Value: objID}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("post not found")
	}

	return nil
}

func (r *repository) CreateComment(ctx context.Context, postID string, comment posts.Comment) (posts.Comment, error) {
	comment.CreatedAt = time.Now().UTC().Round(time.Millisecond)
	comment.UpdateAt = time.Now().UTC().Round(time.Millisecond)
	if comment.UserID == "" {
		return comment, errors.New("user ID is required")
	}
	id, err := uuid.NewV4()
	if err != nil {
		return comment, err
	}
	comment.ID = id.String()

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return comment, err
	}

	if _, err = r.db.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{{Key: "$push", Value: bson.D{
		{Key: "comments", Value: comment},
	}}}); err != nil {
		return comment, err
	}

	return comment, nil
}

func (r *repository) RetrieveCommentByID(ctx context.Context, id string) (posts.Comment, error) {
	opt := options.FindOne().SetProjection(bson.D{{Key: "comments.$", Value: 1}})

	var result posts.Post
	if err := r.db.FindOne(ctx, bson.D{{Key: "comments._id", Value: id}}, opt).Decode(&result); err != nil {
		return posts.Comment{}, err
	}
	if len(result.Comments) == 0 {
		return posts.Comment{}, errors.New("comment not found")
	}

	return result.Comments[0], nil
}

func (r *repository) RetrieveAllComments(ctx context.Context, page posts.Page) (posts.CommentsPage, error) { //nolint:dupl
	if page.PostID == "" {
		return posts.CommentsPage{}, errors.New("post ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(page.PostID)
	if err != nil {
		return posts.CommentsPage{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}

	opts := options.Find()
	opts.SetProjection(bson.D{{Key: "comments", Value: bson.D{{Key: "$slice", Value: bson.A{page.Offset, page.Limit}}}}})
	opts.SetSort(primitive.D{{Key: "comments.created_at", Value: -1}})

	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return posts.CommentsPage{}, err
	}
	defer cursor.Close(ctx)

	var commentsPage posts.CommentsPage
	for cursor.Next(ctx) {
		var post posts.Post
		if err := cursor.Decode(&post); err != nil {
			return commentsPage, err
		}
		commentsPage.Comments = append(commentsPage.Comments, post.Comments...)
	}

	post, err := r.RetrieveByID(ctx, page.PostID)
	if err != nil {
		return posts.CommentsPage{}, err
	}

	commentsPage.Total = uint64(len(post.Comments))
	commentsPage.Offset = page.Offset
	commentsPage.Limit = page.Limit

	return commentsPage, nil
}

func (r *repository) UpdateComment(ctx context.Context, comment posts.Comment) (posts.Comment, error) {
	comment.UpdateAt = time.Now().UTC().Round(time.Millisecond)
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "comments._id", Value: comment.ID}}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "comments.$", Value: comment},
	}}})
	if err != nil {
		return comment, err
	}
	if result.MatchedCount == 0 {
		return comment, errors.New("comment not found")
	}

	return comment, nil
}

func (r *repository) DeleteComment(ctx context.Context, id string) error {
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "comments._id", Value: id}}, bson.D{{Key: "$pull", Value: bson.D{
		{Key: "comments", Value: bson.D{{Key: "_id", Value: id}}},
	}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func (r *repository) CreateLike(ctx context.Context, postID string, like posts.Like) (posts.Like, error) {
	if like.UserID == "" {
		return like, errors.New("user ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return like, err
	}
	like.CreatedAt = time.Now().UTC().Round(time.Millisecond)

	result, err := r.db.UpdateOne(ctx, primitive.M{"_id": objID}, primitive.M{"$push": primitive.M{"likes": like}})
	if err != nil {
		return like, err
	}
	if result.MatchedCount == 0 {
		return like, errors.New("post not found")
	}

	return like, nil
}

func (r *repository) RetrieveAllLikes(ctx context.Context, page posts.Page) (posts.LikesPage, error) { //nolint:dupl
	if page.PostID == "" {
		return posts.LikesPage{}, errors.New("post ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(page.PostID)
	if err != nil {
		return posts.LikesPage{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}

	opts := options.Find()
	opts.SetProjection(bson.D{{Key: "likes", Value: bson.D{{Key: "$slice", Value: bson.A{page.Offset, page.Limit}}}}})
	opts.SetSort(primitive.D{{Key: "likes.created_at", Value: -1}})

	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return posts.LikesPage{}, err
	}
	defer cursor.Close(ctx)

	var likesPage posts.LikesPage
	for cursor.Next(ctx) {
		var post posts.Post
		if err := cursor.Decode(&post); err != nil {
			return likesPage, err
		}
		likesPage.Likes = append(likesPage.Likes, post.Likes...)
	}
	post, err := r.RetrieveByID(ctx, page.PostID)
	if err != nil {
		return posts.LikesPage{}, err
	}
	likesPage.Total = uint64(len(post.Likes))
	likesPage.Offset = page.Offset
	likesPage.Limit = page.Limit

	return likesPage, nil
}

func (r *repository) DeleteLike(ctx context.Context, userID string) error {
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "likes.user_id", Value: userID}}, bson.D{{Key: "$pull", Value: bson.D{
		{Key: "likes", Value: bson.D{{Key: "user_id", Value: userID}}},
	}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *repository) CreateShare(ctx context.Context, postID string, share posts.Share) (posts.Share, error) {
	if share.UserID == "" {
		return share, errors.New("user ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return share, err
	}
	share.CreatedAt = time.Now().UTC().Round(time.Millisecond)

	_, err = r.db.UpdateOne(ctx, primitive.M{"_id": objID}, primitive.M{"$push": primitive.M{"shares": share}})
	if err != nil {
		return share, err
	}

	return share, nil
}

func (r *repository) RetrieveAllShares(ctx context.Context, page posts.Page) (posts.SharesPage, error) { //nolint:dupl
	if page.PostID == "" {
		return posts.SharesPage{}, errors.New("post ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(page.PostID)
	if err != nil {
		return posts.SharesPage{}, err
	}
	filter := bson.D{{Key: "_id", Value: objID}}

	opts := options.Find()
	opts.SetProjection(bson.D{{Key: "shares", Value: bson.D{{Key: "$slice", Value: bson.A{page.Offset, page.Limit}}}}})
	opts.SetSort(primitive.D{{Key: "shares.created_at", Value: -1}})

	cursor, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return posts.SharesPage{}, err
	}
	defer cursor.Close(ctx)

	var sharesPage posts.SharesPage
	for cursor.Next(ctx) {
		var post posts.Post
		if err := cursor.Decode(&post); err != nil {
			return sharesPage, err
		}
		sharesPage.Shares = append(sharesPage.Shares, post.Shares...)
	}

	post, err := r.RetrieveByID(ctx, page.PostID)
	if err != nil {
		return posts.SharesPage{}, err
	}
	sharesPage.Total = uint64(len(post.Shares))
	sharesPage.Offset = page.Offset
	sharesPage.Limit = page.Limit

	return sharesPage, nil
}

func (r *repository) DeleteShare(ctx context.Context, userID string) error {
	result, err := r.db.UpdateOne(ctx, bson.D{{Key: "shares.user_id", Value: userID}}, bson.D{{Key: "$pull", Value: bson.D{
		{Key: "shares", Value: bson.D{{Key: "user_id", Value: userID}}},
	}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("share not found")
	}

	return nil
}

func getFilter(page posts.Page) (bson.D, *options.FindOptions) {
	filter := bson.D{}
	opts := options.Find()

	opts.SetSort(primitive.D{{Key: "created_at", Value: -1}}).SetLimit(int64(page.Limit)).SetSkip(int64(page.Offset))

	if page.Tag != "" {
		filter = append(filter, bson.E{Key: "tags", Value: bson.D{{Key: "$all", Value: bson.A{page.Tag}}}})
	}
	if page.Visibility {
		filter = append(filter, bson.E{Key: "visibility", Value: true})
	}
	if page.PostID != "" {
		objID, err := primitive.ObjectIDFromHex(page.PostID)
		if err != nil {
			return bson.D{}, opts
		}
		filter = append(filter, bson.E{Key: "_id", Value: objID})
	}

	if page.UserID != "" {
		filter = append(filter, bson.E{Key: "user_id", Value: page.UserID})
	}

	return filter, opts
}
