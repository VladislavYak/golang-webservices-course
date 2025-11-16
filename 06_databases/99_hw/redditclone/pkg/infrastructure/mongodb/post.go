package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ post.PostRepository = new(PostRepoMongo)

type PostRepoMongo struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

func NewPostRepoMongo(client *mongo.Client, dbName, collectionName string) *PostRepoMongo {

	collection := client.Database(dbName).Collection(collectionName)

	return &PostRepoMongo{
		Collection: collection,
		Client:     client,
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
	const op = "GetAllPosts"
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var postsTmp []*postTmp

	if err = cursor.All(context.TODO(), &postsTmp); err != nil {
		return nil, errors.Wrap(err, op)
	}

	var Posts []*post.Post
	for _, postIter := range postsTmp {

		Posts = append(Posts, postIter.ToPost())

	}

	return Posts, nil
}

func (pp *PostRepoMongo) GetPostsByCategoryName(ctx context.Context, CategoryName string) ([]*post.Post, error) {
	const op = "GetPostsByCategoryName"
	cursor, err := pp.Collection.Find(context.TODO(), bson.D{{Key: "category", Value: CategoryName}})

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var postsTmp []*postTmp
	if err = cursor.All(context.TODO(), &postsTmp); err != nil {
		return nil, errors.Wrap(err, op)
	}

	var Posts []*post.Post
	for _, postIter := range postsTmp {

		Posts = append(Posts, postIter.ToPost())

	}
	return Posts, nil
}

func (pp *PostRepoMongo) GetPostByID(ctx context.Context, ID string) (*post.Post, error) {
	const op = "GetPostByID"
	value, _ := bson.ObjectIDFromHex(ID)

	filter := bson.M{"_id": value}

	var postTmp *postTmp
	err := pp.Collection.FindOne(context.TODO(), filter).Decode(&postTmp)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	Post := postTmp.ToPost()

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

func (pp *PostRepoMongo) AddPost(ctx context.Context, Post *post.Post) (*post.Post, error) {
	const op = "AddPost"

	// prettify it somehow
	p := &postTmp{
		// Id:               pt.ObjectID.Hex(),
		Category:         Post.Category,
		Type:             Post.Type,
		Url:              Post.Url,
		Text:             Post.Text,
		Title:            Post.Title,
		Votes:            Post.Votes,
		Comments:         Post.Comments,
		Created:          Post.Created,
		UpvotePercentage: Post.UpvotePercentage,
		Score:            Post.Score,
		Views:            Post.Views,
		Author:           Post.Author,
	}

	result, err := pp.Collection.InsertOne(context.TODO(), p)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	ReturnedPost := p.ToPost()

	if objID, ok := result.InsertedID.(bson.ObjectID); ok {
		// Преобразуем ObjectID в строку (hex-формат)
		IdConverted := objID.Hex()
		ReturnedPost.Id = IdConverted
	} else {
		return nil, errors.New("cannot convert id")
	}

	return ReturnedPost, nil
}

func (pp *PostRepoMongo) DeletePost(ctx context.Context, Id string) (*post.Post, error) {
	const op = "DeletePost"

	value, err := bson.ObjectIDFromHex(Id)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	_, err = pp.Collection.DeleteOne(context.TODO(), bson.D{{"_id", value}})
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	// yakovlev: check this out
	return nil, nil

	// return nil, errors.New("this id doesnot exist")

}

// func (pp *PostRepoMongo) Upvote(ctx context.Context, PostID string) (*post.Post, error) {
// 	const op = "Upvote"
// 	userID, ok := ctx.Value("UserID").(string)
// 	if !ok {
// 		return nil, errors.New("cannot cast userID to string")
// 	}
// 	// Convert postID to ObjectID
// 	objID, err := bson.ObjectIDFromHex(PostID)
// 	if err != nil {
// 		return nil, errors.Wrap(err, op)
// 	}
// 	filter := bson.M{
// 		"_id": objID,
// 	}

// 	// Aggregation pipeline для обновления votes
// 	update := bson.A{
// 		bson.M{
// 			"$set": bson.M{
// 				"votes": bson.M{
// 					"$cond": bson.M{
// 						"if": bson.M{
// 							"$in": []interface{}{userID, "$votes.user"},
// 						},
// 						"then": bson.M{
// 							"$map": bson.M{
// 								"input": "$votes",
// 								"as":    "vote",
// 								"in": bson.M{
// 									"$cond": bson.M{
// 										"if": bson.M{"$eq": []interface{}{"$$vote.user", userID}},
// 										"then": bson.M{
// 											"user": "$$vote.user",
// 											"vote": 1,
// 										},
// 										"else": "$$vote",
// 									},
// 								},
// 							},
// 						},
// 						"else": bson.M{
// 							"$concatArrays": []interface{}{
// 								"$votes",
// 								[]bson.M{{
// 									"user": userID,
// 									"vote": 1,
// 								}},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Опции для возврата обновлённого документа
// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

// 	// Выполнение атомарного обновления
// 	var tmpPost postTmp
// 	err = pp.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&tmpPost)

// 	ReturnedPost := tmpPost.ToPost()

// 	if err == mongo.ErrNoDocuments {
// 		return nil, domain.PostNotFoundError
// 	}
// 	if err != nil {
// 		return nil, errors.Wrap(err, op)
// 	}

// 	return ReturnedPost, nil
// }

// // yakovlev: this is almost full copy of Upvote, but lazy now
// func (pp *PostRepoMongo) Downvote(ctx context.Context, PostID string) (*post.Post, error) {
// 	const op = "Downvote"
// 	userID, ok := ctx.Value("UserID").(string)
// 	if !ok {
// 		return nil, errors.New("cannot cast userID to string")
// 	}

// 	// Convert postID to ObjectID
// 	objID, err := bson.ObjectIDFromHex(PostID)
// 	if err != nil {
// 		return nil, errors.Wrap(err, op)
// 	}
// 	filter := bson.M{
// 		"_id": objID,
// 	}

// 	// Aggregation pipeline для обновления votes
// 	update := bson.A{
// 		bson.M{
// 			"$set": bson.M{
// 				"votes": bson.M{
// 					"$cond": bson.M{
// 						"if": bson.M{
// 							"$in": []interface{}{userID, "$votes.user"},
// 						},
// 						"then": bson.M{
// 							"$map": bson.M{
// 								"input": "$votes",
// 								"as":    "vote",
// 								"in": bson.M{
// 									"$cond": bson.M{
// 										"if": bson.M{"$eq": []interface{}{"$$vote.user", userID}},
// 										"then": bson.M{
// 											"user": "$$vote.user",
// 											"vote": -1,
// 										},
// 										"else": "$$vote",
// 									},
// 								},
// 							},
// 						},
// 						"else": bson.M{
// 							"$concatArrays": []interface{}{
// 								"$votes",
// 								[]bson.M{{
// 									"user": userID,
// 									"vote": 1,
// 								}},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Опции для возврата обновлённого документа
// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

// 	// Выполнение атомарного обновления
// 	var tmpPost postTmp
// 	err = pp.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&tmpPost)

// 	if err == mongo.ErrNoDocuments {
// 		return nil, fmt.Errorf("post with ID %s not found", PostID)
// 	}
// 	if err != nil {
// 		return nil, errors.Wrap(err, op)
// 	}

// 	ReturnedPost := tmpPost.ToPost()

// 	return ReturnedPost, nil

// }

func (pp *PostRepoMongo) Vote(ctx context.Context, PostID string, vote int) (*post.Post, error) {
	const op = "Vote"
	if vote != 1 && vote != -1 {
		return nil, errors.New("vote must be 1 (upvote) or -1 (downvote)")
	}

	userID, ok := ctx.Value("UserID").(string)
	if !ok {
		return nil, errors.New("cannot cast userID to string")
	}

	// Convert postID to ObjectID
	objID, err := bson.ObjectIDFromHex(PostID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	filter := bson.M{"_id": objID}

	// Универсальный pipeline
	update := bson.A{
		bson.M{
			"$set": bson.M{
				"votes": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$in": []interface{}{userID, "$votes.user"},
						},
						"then": bson.M{
							"$map": bson.M{
								"input": "$votes",
								"as":    "v",
								"in": bson.M{
									"$cond": bson.M{
										"if":   bson.M{"$eq": []interface{}{"$$v.user", userID}},
										"then": bson.M{"user": "$$v.user", "vote": vote}, // ← вот и всё!
										"else": "$$v",
									},
								},
							},
						},
						"else": bson.M{
							"$concatArrays": []interface{}{
								"$votes",
								[]bson.M{{"user": userID, "vote": vote}},
							},
						},
					},
				},
			},
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var tmpPost postTmp
	err = pp.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&tmpPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, post.PostNotFoundError
		}
		return nil, errors.Wrap(err, op)
	}

	return tmpPost.ToPost(), nil
}

func (pp *PostRepoMongo) Unvote(ctx context.Context, PostId string) (*post.Post, error) {
	const op = "Unvote"
	userID, ok := ctx.Value("UserID").(string)
	if !ok {
		return nil, errors.New("cannot cast userID to string")
	}

	objID, err := bson.ObjectIDFromHex(PostId)
	if err != nil {
		return nil, post.InvalidPostIdError
	}
	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$pull": bson.M{
			"votes": bson.M{
				"user": userID,
			},
		},
	}

	// Опции для возврата обновлённого документа
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var tmpPost postTmp
	err = pp.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&tmpPost)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("post with ID %s not found", PostId)
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	ReturnedPost := tmpPost.ToPost()

	return ReturnedPost, nil
}

// UpdateScore calculates the sum of votes.vote and updates the score field
func (pp *PostRepoMongo) UpdateScore(ctx context.Context, PostId string) (*post.Post, error) {
	const op = "UpdateScore"

	// Convert postID to ObjectID
	objID, err := bson.ObjectIDFromHex(PostId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	// Fetch the post
	var tmpPost postTmp
	err = pp.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&tmpPost)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("post with ID %s not found", PostId)
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	// Calculate score by summing votes.vote
	updatedScore := 0
	for _, vote := range tmpPost.Votes {
		updatedScore += vote.VoteScore
	}

	// Update the score in MongoDB
	_, err = pp.Collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "score", Value: updatedScore}}},
	})
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	// Update the post object
	tmpPost.Score = updatedScore
	ReturnedPost := tmpPost.ToPost()
	return ReturnedPost, nil
}
