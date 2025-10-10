package main

import (
	"context"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("testing").Collection("numbers")

	// MyPost := post.NewPost("music", "text", "", "asdiojasdasd", "waswawa", *user.NewUser("vlad"))

	// result, _ := collection.InsertOne(context.TODO(), MyPost)

	// fmt.Println("result", result.InsertedID)

	cursor, err := collection.Find(context.TODO(), bson.D{})

	if err != nil {
		panic(err)
	}

	var Posts []*post.Post
	if err := cursor.All(context.TODO(), &Posts); err != nil {
		// change panic to something
		panic(err)
	}
	// fmt.Println("posts", Posts[0].MongoId)

	// fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	// res, _ := collection.InsertOne(ctx, bson.D{{Key: "name", Value: "pi"}, {Key: "value", Value: 3.14159}})
	// id := res.InsertedID

	// fmt.Println("id", id)
}

// https://deepwiki.com/mongodb/mongo-go-driver/4.1-marshalunmarshal
