package main

import (
	"context"
	"fmt"

	"github.com/VladislavYak/redditclone/mocks"
	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"go.uber.org/mock/gomock"
)

func main() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	// mockRepoComment := mocks.NewMockCommentRepository(ctrl)

	newPost := post.NewPost("test", "test", "test", "test", "test", *user.NewUser("Vlad"))

	mockPostRepo.EXPECT().AddPost(gomock.Any(), newPost).Return(post.NewPost("test1", "test1", "test1", "test1", "test1", *user.NewUser("Egor")), nil)

	ctx := context.Background()

	application.NewPostImpl(mockPostRepo)

	post, err := mockPostRepo.AddPost(ctx, newPost)

	fmt.Println("post", post)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS: Created post %+v\n", post)

}
