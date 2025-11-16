package mongodb

import (
	"context"
	"fmt"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ comment.CommentRepository = new(CommentRepoMongo)

type CommentRepoMongo struct {
	Collection *mongo.Collection
}

func NewCommentRepoMongo(client *mongo.Client, dbName string, collectionName string) *CommentRepoMongo {

	collection := client.Database(dbName).Collection(collectionName)

	return &CommentRepoMongo{
		Collection: collection,
	}
}

func (cr *CommentRepoMongo) AddComment(ctx context.Context, Id string, Comment *comment.Comment) error {
	const op = "Add Comment"

	Comment.WithId(bson.NewObjectID().Hex())

	objID, err := bson.ObjectIDFromHex(Id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	update := bson.M{
		"$push": bson.M{
			"comments": Comment,
		},
	}

	result, err := cr.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
	)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if result.MatchedCount == 0 {
		return post.PostNotFoundError
	}

	return nil

}

func (cr *CommentRepoMongo) DeleteComment(ctx context.Context, Id string, CommentId string) error {

	objID, err := bson.ObjectIDFromHex(Id)
	if err != nil {
		return post.InvalidPostIdError
	}

	// Define the update operation using $pull
	update := bson.M{
		"$pull": bson.M{
			"comments": bson.M{
				"id": CommentId,
			},
		},
	}

	// Update the post in MongoDB
	result, err := cr.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if result.MatchedCount == 0 {
		return post.PostNotFoundError
	}
	if result.ModifiedCount == 0 {
		return comment.CommentNotFoundError
	}

	return nil

}
