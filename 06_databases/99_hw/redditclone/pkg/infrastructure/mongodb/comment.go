package mongodb

import (
	"context"
	"fmt"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CommentRepoMongo struct {
	Collection *mongo.Collection
}

func NewCommentRepoMongo(client *mongo.Client, dbName string, collectionName string) *CommentRepoMongo {

	collection := client.Database(dbName).Collection(collectionName)

	return &CommentRepoMongo{
		Collection: collection,
	}
}

func (cr *CommentRepoMongo) AddComment(Id string, Comment *comment.Comment) error {
	fmt.Println("inserteing comment")

	Comment.WithId(bson.NewObjectID().Hex())

	objID, err := bson.ObjectIDFromHex(Id)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	update := bson.M{
		"$push": bson.M{
			"comments": Comment,
		},
	}

	result, err := cr.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		update,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no post found with ID %s", Id)
	}

	return nil

}

func (cr *CommentRepoMongo) DeleteComment(Id string, CommentId string) error {

	objID, err := bson.ObjectIDFromHex(Id)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
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
		context.TODO(),
		bson.M{"_id": objID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no post found with ID %s", Id)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no comment found with ID %s in post %s", CommentId, Id)
	}

	return nil

}
