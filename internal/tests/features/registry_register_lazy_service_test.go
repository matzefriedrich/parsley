package features

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Registry_register_lazy_service_type_factory(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = features.RegisterLazy[*foo](registry, func() *foo {
		return &foo{}
	}, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(context.Background())

	// Act
	lazyFactory, err := resolving.ResolveRequiredService[func() features.Lazy[*foo]](resolver, ctx)
	actual := lazyFactory()
	fooInstance0 := actual.Value()
	fooInstance1 := actual.Value()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, lazyFactory)
	assert.NotNil(t, actual)
	assert.NotNil(t, fooInstance0)
	assert.NotNil(t, fooInstance1)
}

type foo struct {
}
