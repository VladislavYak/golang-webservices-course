package post

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type PostRepo struct {
	Data []*Post
	*sync.Mutex
	lastID    int
	commentID int
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		[]*Post{},
		&sync.Mutex{},
		0,
		0,
	}
}

func (pp *PostRepo) GetAllPosts() ([]*Post, error) {
	return pp.Data, nil
}

func (pp *PostRepo) GetPostsByCategoryName(CategoryName string) ([]*Post, error) {
	res := []*Post{}

	for _, post := range pp.Data {
		if post.Category == CategoryName {
			res = append(res, post)
		}
	}
	return res, nil

}

func (pp *PostRepo) GetPostByID(ID string) (*Post, error) {
	for _, post := range pp.Data {
		if post.Id == ID {
			return post, nil
		}
	}
	return nil, errors.New("not found")
}

func (pp *PostRepo) GetPostsByUsername(Username string) ([]*Post, error) {
	res := []*Post{}

	for _, post := range pp.Data {
		if post.Author.Username == Username {
			res = append(res, post)
		}
	}
	return res, nil

}

func (pp *PostRepo) UpdatePostViews(ID string) error {
	for _, Post := range pp.Data {
		if Post.Id == ID {
			Post.Views += 1
			return nil
		}
	}
	return errors.New("not found")
}

func (pp *PostRepo) AddPost(Post *Post) (*Post, error) {
	// pretty random handling mutexes
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	Post = Post.WithId(strconv.Itoa(pp.lastID))

	pp.lastID++

	pp.Data = append(pp.Data, Post)

	fmt.Println("my Posts", pp.Data)
	return Post, nil
}

func (pp *PostRepo) DeletePost(Id string) (*Post, error) {
	for i, value := range pp.Data {
		if value.Id == Id {
			pp.Data = append(pp.Data[:i], pp.Data[i+1:]...)
		}
		return value, nil
	}

	return nil, errors.New("this id doesnot exist")

}

// func (pp *PostRepo) AddComment(Id string, comment *comment.Comment) (*Post, error) {
// 	// add more mutexes handling
// 	pp.Mutex.Lock()
// 	defer pp.Mutex.Unlock()

// 	for _, Post := range pp.Data {
// 		if Post.Id == Id {
// 			Post.Comments = append(Post.Comments, *comment.WithId(strconv.Itoa(pp.commentID)))

// 			pp.commentID++
// 			return Post, nil
// 		}
// 	}

// 	return nil, errors.New("post not found")
// }

// func (pp *PostRepo) DeleteComment(id string, commentId string) (*Post, error) {

// 	pp.Mutex.Lock()
// 	defer pp.Mutex.Unlock()
// 	for i, post := range pp.Data {
// 		if post.Id == id {

// 			for j, comment := range post.Comments {
// 				if comment.Id == commentId {
// 					post.Comments = append(post.Comments[:j], post.Comments[j+1:]...)
// 					pp.Data[i] = post
// 					return post, nil
// 				}

// 			}

// 		}

// 	}
// 	return nil, errors.New("this id doesnot exist")
// }

// yakovlev: add proper error handling
func (pp *PostRepo) Upvote(id string, user_id string) (*Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == id {
			for j, voteIter := range Post.Votes {
				if voteIter.User == user_id {

					pp.Data[i].Votes[j].WithVote(1)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			pp.Data[i].Votes = append(pp.Data[i].Votes, Vote{User: user_id, VoteScore: 1})
			// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

			pp.Data[i].UpdateScore()

			return pp.Data[i], nil
		}
	}

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepo) Downvote(id string, user_id string) (*Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == id {
			for j, voteIter := range Post.Votes {
				if voteIter.User == user_id {

					pp.Data[i].Votes[j].WithVote(-1)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			pp.Data[i].Votes = append(pp.Data[i].Votes, Vote{User: user_id, VoteScore: -1})
			// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

			pp.Data[i].UpdateScore()

			return pp.Data[i], nil
		}
	}

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepo) Unvote(id string, user_id string) (*Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == id {
			for j, voteIter := range Post.Votes {
				if voteIter.User == user_id {

					pp.Data[i].Votes = append(pp.Data[i].Votes[:j], pp.Data[i].Votes[j+1:]...)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			return nil, errors.New("cannot find user for specified post")
		}
	}

	return nil, errors.New("this id doesnot exist")
}
