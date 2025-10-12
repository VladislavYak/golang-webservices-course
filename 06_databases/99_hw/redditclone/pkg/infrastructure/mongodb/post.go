package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ post.PostRepository = new(PostRepoMongo)

type PostRepoMongo struct {
	Collection *mongo.Collection
}

func NewPostRepoMongo(client *mongo.Client, dbName, collectionName string) *PostRepoMongo {

	collection := client.Database(dbName).Collection(collectionName)

	return &PostRepoMongo{
		Collection: collection,
	}
}

type postTmp struct {
	ObjectID         bson.ObjectID     `bson:"_id,omitempty"`
	Category         string            `json:"category"`
	Type             string            `json:"type"`
	Url              string            `json:"url,omitempty"`
	Text             string            `json:"text,omitempty"`
	Title            string            `json:"title"`
	Votes            []post.Vote       `json:"votes"`
	Comments         []comment.Comment `json:"comments"`
	Created          time.Time         `json:"created"`
	UpvotePercentage int               `json:"upvotePercentage"`
	Score            int               `json:"score"`
	Views            int               `json:"views"`
	Author           user.User         `json:"author"`
}

func (pt *postTmp) ToPost() *post.Post {
	return &post.Post{
		Id:               pt.ObjectID.Hex(),
		Category:         pt.Category,
		Type:             pt.Type,
		Url:              pt.Url,
		Text:             pt.Text,
		Title:            pt.Title,
		Votes:            pt.Votes,
		Comments:         pt.Comments,
		Created:          pt.Created,
		UpvotePercentage: pt.UpvotePercentage,
		Score:            pt.Score,
		Views:            pt.Views,
		Author:           pt.Author,
	}
}

func (pp *PostRepoMongo) GetAllPosts(ctx context.Context) ([]*post.Post, error) {
	fmt.Println("inside GetAllPosts")
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println("я в курсоре")
		panic(err)
	}

	var postsTmp []*postTmp

	if err = cursor.All(context.TODO(), &postsTmp); err != nil {
		panic(err)
	}

	var Posts []*post.Post
	for _, postIter := range postsTmp {

		Posts = append(Posts, postIter.ToPost())

	}

	return Posts, nil
}

func (pp *PostRepoMongo) GetPostsByCategoryName(ctx context.Context, CategoryName string) ([]*post.Post, error) {
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{{Key: "category", Value: CategoryName}})

	if err != nil {
		panic(err)
	}
	// var Posts []*post.Post
	var postsTmp []*postTmp
	if err = cursor.All(context.TODO(), &postsTmp); err != nil {
		panic(err)
	}

	var Posts []*post.Post
	for _, postIter := range postsTmp {

		Posts = append(Posts, postIter.ToPost())

	}
	return Posts, nil
}

func (pp *PostRepoMongo) GetPostByID(ctx context.Context, ID string) (*post.Post, error) {
	value, _ := bson.ObjectIDFromHex(ID)

	fmt.Println("GetPostByID, value", value)

	filter := bson.M{"_id": value}

	var postTmp *postTmp
	err := pp.Collection.FindOne(context.TODO(), filter).Decode(&postTmp)

	Post := postTmp.ToPost()

	if err != nil {
		panic(err)
	}

	fmt.Println("Post GetPostByID", Post)

	return Post, nil
}

func (pp *PostRepoMongo) GetPostsByUsername(ctx context.Context, Username string) ([]*post.Post, error) {
	cursor, err := pp.Collection.Find(context.TODO(), bson.M{"author.username": Username})

	if err != nil {
		panic(err)
	}

	var postsTmp []*postTmp
	if err = cursor.All(context.TODO(), &postsTmp); err != nil {
		panic(err)
	}

	var Posts []*post.Post
	for _, postIter := range postsTmp {

		Posts = append(Posts, postIter.ToPost())

	}

	return Posts, nil
}

func (pp *PostRepoMongo) UpdatePostViews(ID string) error {

	value, _ := bson.ObjectIDFromHex(ID)

	fmt.Println("GetPostByID, value", value)

	filter := bson.D{{Key: "_id", Value: value}}
	update := bson.D{{"$inc", bson.M{"views": 1}}}

	result, err := pp.Collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Documents matched: %v\n", result.MatchedCount)
		fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
		return nil
	}

	// return &Post, nil
	return errors.New("not found")
}

func (pp *PostRepoMongo) AddPost(ctx context.Context, Post *post.Post) (*post.Post, error) {

	result, _ := pp.Collection.InsertOne(context.TODO(), Post)

	fmt.Println("inserted id", result.InsertedID)
	return Post, nil
}

func (pp *PostRepoMongo) DeletePost(ctx context.Context, Id string) (*post.Post, error) {

	value, _ := bson.ObjectIDFromHex(Id)

	_, err := pp.Collection.DeleteOne(context.TODO(), bson.D{{"_id", value}})
	if err != nil {
		return nil, err
	}

	return nil, errors.New("this id doesnot exist")

}

// yakovlev: add proper error handling
func (pp *PostRepoMongo) Upvote(ctx context.Context, id string, UserId string) (*post.Post, error) {

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepoMongo) Downvote(ctx context.Context, id string, UserId string) (*post.Post, error) {

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepoMongo) Unvote(ctx context.Context, id string, UserId string) (*post.Post, error) {

	return nil, errors.New("this id doesnot exist")
}
