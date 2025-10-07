package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CommentRepoMongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewCommentRepoMongo() *CommentRepoMongo {

	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	collection := client.Database("testing").Collection("posts")

	err := client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection successfullt initialized")
	return &CommentRepoMongo{
		client, collection,
	}
}
