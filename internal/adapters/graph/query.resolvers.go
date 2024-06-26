package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"post-and-comments/internal/adapters/graph/generated"
	"post-and-comments/internal/adapters/graph/grmodels"
	"post-and-comments/internal/models"
)

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, postLimit int, postOffset int, commentLimit int, commentOffset int, depthLimit int) ([]*grmodels.Post, error) {
	posts, err := r.PostSrv.Posts(ctx, postLimit, postOffset)
	if err != nil {
		return nil, err
	}

	grposts := make([]*grmodels.Post, 0, len(posts))

	for _, post := range posts {
		comments, err := r.CommentSrv.Comments(ctx, post.ID, commentOffset, commentLimit, depthLimit)
		if err != nil {
			return nil, err
		}

		grposts = append(grposts, r.toGrPosts(ctx, post, comments))
	}

	return grposts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, postID string, commentLimit int, commentOffset int, depthLimit int, afterComment *string) (*grmodels.Post, error) {
	mPostID, err := models.ParseID(postID)
	if err != nil {
		return nil, err
	}

	post, err := r.PostSrv.Post(ctx, mPostID)
	if err != nil {
		return nil, err
	}

	var comments []*models.Comment
	if afterComment == nil {
		comments, err = r.CommentSrv.Comments(ctx, post.ID, commentOffset, commentLimit, depthLimit)
		if err != nil {
			return nil, err
		}
	} else {
		id, err := models.ParseID(*afterComment)
		if err != nil {
			return nil, err
		}

		comments, err = r.CommentSrv.CommentsAfter(ctx, mPostID, id, commentOffset, commentLimit, depthLimit)
	}

	return r.toGrPosts(ctx, post, comments), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
