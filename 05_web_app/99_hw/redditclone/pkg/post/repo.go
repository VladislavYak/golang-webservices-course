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

func (pp *PostRepo) AddComment(Id string, comment *Comment) (*Post, error) {
	// add more mutexes handling
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for _, Post := range pp.Data {
		if Post.Id == Id {
			Post.Comments = append(Post.Comments, *comment.WithId(strconv.Itoa(pp.commentID)))

			pp.commentID++
			return Post, nil
		}
	}

	return nil, errors.New("post not found")
}

// yakovlev move DeleteComment here
func (pp *PostRepo) DeleteComment(id string, commentId string) (*Post, error) {

	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()
	for i, post := range pp.Data {
		if post.Id == id {

			for j, comment := range post.Comments {
				if comment.Id == commentId {
					post.Comments = append(post.Comments[:j], post.Comments[j+1:]...)
					pp.Data[i] = post
					return post, nil
				}

			}

		}

	}
	return nil, errors.New("this id doesnot exist")
}
