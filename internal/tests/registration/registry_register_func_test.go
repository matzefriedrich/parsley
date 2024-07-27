package registration

import (
	"context"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_RegisterTransient_registers_factory_function_to_resolve_dynamic_dependency_at_runtime(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	_ = registration.RegisterTransient(sut, newPersonFactory)
	_ = registration.RegisterTransient(sut, newUserController)

	r := resolving.NewResolver(sut)

	const userId = "123"

	// Act
	controller, _ := resolving.ResolveRequiredService[userController](r, context.Background())
	model := controller.GetUserInfo(userId)

	userServiceFactory, _ := resolving.ResolveRequiredService[func(string) (userService, error)](r, context.Background())
	service, factoryErr := userServiceFactory(userId)

	// Assert
	assert.NoError(t, factoryErr)
	assert.NotNil(t, service)

	assert.Equal(t, userId, model.Id)
	assert.Equal(t, userId, service.UserId())
}

type idUserService struct {
	id string
}

func (f *idUserService) UserId() string {
	return f.id
}

type userModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (f *idUserService) UserName() string {
	return "" // retrieve username from data backend
}

type idUserController struct {
	userServiceFactory func(name string) (userService, error)
}

func (f *idUserController) GetUserInfo(userId string) userModel {
	userService, _ := f.userServiceFactory(userId)
	return userModel{
		Id:   userId,
		Name: userService.UserName(),
	}
}

type userService interface {
	UserName() string
	UserId() string
}

type userController interface {
	GetUserInfo(userId string) userModel
}

func newPersonFactory() func(string) (userService, error) {
	return func(userId string) (userService, error) {
		return &idUserService{
			id: userId,
		}, nil
	}
}

func newUserController(userServiceFactory func(string) (userService, error)) userController {
	return &idUserController{
		userServiceFactory: userServiceFactory,
	}
}
