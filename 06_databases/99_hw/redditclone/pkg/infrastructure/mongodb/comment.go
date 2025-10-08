package mongodb

import (
	"github.com/VladislavYak/redditclone/pkg/domain/comment"
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

	return nil

}

func (cr *CommentRepoMongo) DeleteComment(Id string, CommentId string) error {

	return nil
}
