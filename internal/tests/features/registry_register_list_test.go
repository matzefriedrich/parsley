package features

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Resolver_register_list_resolver(t *testing.T) {

	// Arrange
	ctx := t.Context()

	registry := registration.NewServiceRegistry()
	registry.Register(newLocalDataService, types.LifetimeTransient)
	registry.Register(newRemoteDataService, types.LifetimeTransient)
	features.RegisterList[dataService](ctx, registry)

	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(ctx)

	// Act
	actual, err := resolving.ResolveRequiredService[[]dataService](resolverContext, resolver)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, actual, 2)
}

func Test_Resolver_resolve_multiple_instances_of_type(t *testing.T) {

	// Arrange
	ctx := t.Context()

	registry := registration.NewServiceRegistry()
	registry.Register(newLocalDataService, types.LifetimeTransient)
	registry.Register(newRemoteDataService, types.LifetimeTransient)
	features.RegisterList[dataService](ctx, registry)
	registry.Register(newControllerWithServiceList, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(ctx)

	// Act
	actual, err := resolving.ResolveRequiredService[*controllerWithServiceList](resolverContext, resolver)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

}

type controllerWithServiceList struct {
	dataServices []dataService
}

func newControllerWithServiceList(dataServices []dataService) *controllerWithServiceList {
	return &controllerWithServiceList{
		dataServices: dataServices,
	}
}
