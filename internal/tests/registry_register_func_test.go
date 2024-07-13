package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Registry_Register_func(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	_ = registration.RegisterTransient(sut, NewPersonFactory)
	_ = registration.RegisterTransient(sut, NewUserController)

	r := resolving.NewResolver(sut)

	const userId = "123"

	// Act
	controller, _ := resolving.ResolveRequiredService[UserController](r, context.Background())
	model := controller.GetUserInfo(userId)

	userServiceFactory, _ := resolving.ResolveRequiredService[func(string) (UserService, error)](r, context.Background())
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

type userModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (f *userService) UserName() string {
	return "" // retrieve username from data backend
}

type userController struct {
	userServiceFactory func(name string) (UserService, error)
}

func (f *userController) GetUserInfo(userId string) userModel {
	userService, _ := f.userServiceFactory(userId)
	return userModel{
		Id:   userId,
		Name: userService.UserName(),
	}
}

type UserService interface {
	UserName() string
	UserId() string
}

type UserController interface {
	GetUserInfo(userId string) userModel
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
