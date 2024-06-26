package graph

import (
	"context"

	"post-and-comments/internal/adapters/graph/grmodels"
	"post-and-comments/internal/models"
	outports "post-and-comments/internal/out-ports"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostSrv    outports.PostService
	CommentSrv outports.CommentService
	UserSrv    outports.UserService
}

func (r *Resolver) toGrComment(ctx context.Context, comment *models.Comment) *grmodels.Comment {
	if comment == nil {
		return nil
	}

	replies := make([]*grmodels.Comment, 0, len(comment.Replies))

	for _, rep := range comment.Replies {
		r.toGrComment(ctx, rep)
	}

	user, err := r.UserSrv.User(ctx, comment.AuthorID)
	if err != nil {
		return nil
	}

	return &grmodels.Comment{
		ID:        comment.ID.String(),
		Author:    toGrUser(user),
		Published: comment.Published,
		Text:      string(comment.Text),
		PostID:    comment.PostID.String(),
		Replies:   replies,
	}
}

func toGrUser(user *models.User) *grmodels.User {
	return &grmodels.User{ID: user.ID.String()}
}

func (r *Resolver) toGrPosts(ctx context.Context, post *models.Post, comments []*models.Comment) *grmodels.Post {
	grcomments := make([]*grmodels.Comment, 0, len(comments))
	for _, comment := range comments {
		grcomments = append(grcomments, r.toGrComment(ctx, comment))
	}

	user, err := r.UserSrv.User(ctx, post.AuthorID)
	if err != nil {
		return nil
	}

	return &grmodels.Post{
		ID:              post.ID.String(),
		Author:          toGrUser(user),
		Text:            string(post.Text),
		Published:       post.Published,
		DisableComments: post.DisableComments,
		Comments:        grcomments,
	}
}
