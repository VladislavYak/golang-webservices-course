package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/VladislavYak/redditclone/mocks"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	fmt.Println("привет")
	ctrl := gomock.NewController(t)
	UserRepo := mocks.NewMockUserRepository(ctrl)
	User := user.NewUser("Vlad")

	ui := UserImpl{ur: UserRepo}

	ctx := context.Background()

	UserRepo.EXPECT().GetUser(ctx, User).Return(&user.User{Username: "Vlad", UserID: "boobar"}, nil)

	_, err := ui.Register(ctx, "Vlad", "boobar")

	fmt.Println("err", err)
	// if err != user.UserAlreadyExistsError {
	// 	t.Error
	// }

}
