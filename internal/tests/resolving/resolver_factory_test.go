package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_Resolver_ResolveRequiredService_factory_function_receives_current_resolver(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(sut, NewFactory)

	r := resolving.NewResolver(sut)

	ctx := resolving.NewScopedContext(context.Background())

	// Act
	serviceFactory, _ := resolving.ResolveRequiredService[FactoryService](r, ctx)
	f := serviceFactory.(*factory)
	actual := f.resolver

	// Assert
	assert.NotNil(t, serviceFactory)

	assert.NotNil(t, actual)
	assert.Equal(t, reflect.ValueOf(r).Pointer(), reflect.ValueOf(actual).Pointer())
}

type factory struct {
	resolver types.Resolver
}

type FactoryService interface {
}

func NewFactory(resolver types.Resolver) FactoryService {
	return &factory{resolver: resolver}
}
