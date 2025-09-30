package post

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PostRepoMongo struct {
	lastID    int
	commentID int

	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewPostRepoMongo() *PostRepoMongo {

	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

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
	return &PostRepoMongo{
		0, 0, client, collection,
	}
}

func (pp *PostRepoMongo) GetAllPosts() ([]*Post, error) {
	fmt.Println("inside GetAllPosts")
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{})

	if err != nil {
		fmt.Println("я в курсоре")
		panic(err)
	}

	var Posts []*Post
	if err = cursor.All(context.TODO(), &Posts); err != nil {
		panic(err)
	}

	for _, post := range Posts {
		post.Id = string(post.MongoId)
	}

	return Posts, nil
}

func (pp *PostRepoMongo) GetPostsByCategoryName(CategoryName string) ([]*Post, error) {
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{{Key: "category", Value: CategoryName}})

	if err != nil {
		panic(err)
	}
	var Posts []*Post
	if err = cursor.All(context.TODO(), &Posts); err != nil {
		panic(err)
	}
	return Posts, nil
}

func (pp *PostRepoMongo) GetPostByID(ID string) (*Post, error) {

	value, _ := bson.ObjectIDFromHex(ID)

	fmt.Println("GetPostByID, value", value)

	filter := bson.M{"_id": value}

	var Post Post
	err := pp.Collection.FindOne(context.TODO(), filter).Decode(&Post)

	if err != nil {
		panic(err)
	}

	fmt.Println("Post GetPostByID", Post)

	return &Post, nil
}

func (pp *PostRepoMongo) GetPostsByUsername(Username string) ([]*Post, error) {
	cursor, err := pp.Collection.Find(context.TODO(), bson.M{"author.username": Username})

	if err != nil {
		panic(err)
	}
	var Posts []*Post
	if err = cursor.All(context.TODO(), &Posts); err != nil {
		panic(err)
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

func (pp *PostRepoMongo) AddPost(Post *Post) (*Post, error) {

	result, _ := pp.Collection.InsertOne(context.TODO(), Post)

	fmt.Println("inserted id", result.InsertedID)
	return Post, nil
}

func (pp *PostRepoMongo) DeletePost(Id string) (*Post, error) {
	// for i, value := range pp.Data {
	// 	if value.Id == Id {
	// 		pp.Data = append(pp.Data[:i], pp.Data[i+1:]...)
	// 	}
	// 	return value, nil
	// }

	return nil, errors.New("this id doesnot exist")

}

func (pp *PostRepoMongo) AddComment(Id string, comment *Comment) (*Post, error) {
	// add more mutexes handling
	// pp.Mutex.Lock()
	// defer pp.Mutex.Unlock()

	// for _, Post := range pp.Data {
	// 	if Post.Id == Id {
	// 		Post.Comments = append(Post.Comments, *comment.WithId(strconv.Itoa(pp.commentID)))

	// 		pp.commentID++
	// 		return Post, nil
	// 	}
	// }

	return nil, errors.New("post not found")
}

func (pp *PostRepoMongo) DeleteComment(id string, commentId string) (*Post, error) {

	// pp.Mutex.Lock()
	// defer pp.Mutex.Unlock()
	// for i, post := range pp.Data {
	// 	if post.Id == id {

	// 		for j, comment := range post.Comments {
	// 			if comment.Id == commentId {
	// 				post.Comments = append(post.Comments[:j], post.Comments[j+1:]...)
	// 				pp.Data[i] = post
	// 				return post, nil
	// 			}

	// 		}

	// 	}

	// }
	return nil, errors.New("this id doesnot exist")
}

// yakovlev: add proper error handling
func (pp *PostRepoMongo) Upvote(id string, user_id string) (*Post, error) {
	// pp.Mutex.Lock()
	// defer pp.Mutex.Unlock()

	// for i, Post := range pp.Data {
	// 	if Post.Id == id {
	// 		for j, voteIter := range Post.Votes {
	// 			if voteIter.User == user_id {

	// 				pp.Data[i].Votes[j].WithVote(1)
	// 				pp.Data[i].UpdateScore()
	// 				return pp.Data[i], nil
	// 			}
	// 		}

	// 		pp.Data[i].Votes = append(pp.Data[i].Votes, Vote{User: user_id, VoteScore: 1})
	// 		// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

	// 		pp.Data[i].UpdateScore()

	// 		return pp.Data[i], nil
	// 	}
	// }

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepoMongo) Downvote(id string, user_id string) (*Post, error) {
	// pp.Mutex.Lock()
	// defer pp.Mutex.Unlock()

	// for i, Post := range pp.Data {
	// 	if Post.Id == id {
	// 		for j, voteIter := range Post.Votes {
	// 			if voteIter.User == user_id {

	// 				pp.Data[i].Votes[j].WithVote(-1)
	// 				pp.Data[i].UpdateScore()
	// 				return pp.Data[i], nil
	// 			}
	// 		}

	// 		pp.Data[i].Votes = append(pp.Data[i].Votes, Vote{User: user_id, VoteScore: -1})
	// 		// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

	// 		pp.Data[i].UpdateScore()

	// 		return pp.Data[i], nil
	// 	}
	// }

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepoMongo) Unvote(id string, user_id string) (*Post, error) {
	// pp.Mutex.Lock()
	// defer pp.Mutex.Unlock()

	// for i, Post := range pp.Data {
	// 	if Post.Id == id {
	// 		for j, voteIter := range Post.Votes {
	// 			if voteIter.User == user_id {

	// 				pp.Data[i].Votes = append(pp.Data[i].Votes[:j], pp.Data[i].Votes[j+1:]...)
	// 				pp.Data[i].UpdateScore()
	// 				return pp.Data[i], nil
	// 			}
	// 		}

	// 		return nil, errors.New("cannot find user for specified post")
	// 	}
	// }

	return nil, errors.New("this id doesnot exist")
}
