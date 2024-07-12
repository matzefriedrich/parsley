package pkg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Registry_Register_func(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	_ = RegisterTransient(sut, NewPersonFactory)
	_ = RegisterTransient(sut, NewUserController)

	r := sut.BuildResolver()

	const userId = "123"

	// Act
	controller, _ := ResolveRequiredService[UserController](r, context.Background())
	model := controller.GetUserInfo(userId)

	userServiceFactory, _ := ResolveRequiredService[func(string) (UserService, error)](r, context.Background())
	service, factoryErr := userServiceFactory(userId)

	// Assert
	assert.NoError(t, factoryErr)
	assert.NotNil(t, service)

	assert.Equal(t, userId, model.Id)
	assert.Equal(t, userId, service.UserId())
}

type userService struct {
	id string
}

func (f *userService) UserId() string {
	return f.id
}

type UserModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (f *userService) UserName() string {
	return "" // retrieve username from data backend
}

type userController struct {
	userServiceFactory func(name string) (UserService, error)
}

func (f *userController) GetUserInfo(userId string) UserModel {
	userService, _ := f.userServiceFactory(userId)
	return UserModel{
		Id:   userId,
		Name: userService.UserName(),
	}
}

type UserService interface {
	UserName() string
	UserId() string
}

type UserController interface {
	GetUserInfo(userId string) UserModel
}

func NewPersonFactory() func(string) (UserService, error) {
	return func(userId string) (UserService, error) {
		return &userService{
			id: userId,
		}, nil
	}
}

func NewUserController(userServiceFactory func(string) (UserService, error)) UserController {
	return &userController{
		userServiceFactory: userServiceFactory,
	}
}
