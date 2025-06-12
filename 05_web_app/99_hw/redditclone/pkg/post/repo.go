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
	lastID int
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		[]*Post{},
		&sync.Mutex{},
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

func (pp *PostRepo) AddPost(Post *Post) (*Post, error) {
	// pretty random handling mutexes
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	Post.Id = strconv.Itoa(pp.lastID)
	pp.lastID++

	pp.Data = append(pp.Data, Post)

	fmt.Println("my Posts", pp.Data)
	return Post, nil
}
