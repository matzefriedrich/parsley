package features

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_register_named_service_resolve_factory(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	err := features.RegisterNamed[dataService](t.Context(), registry,
		registration.NamedServiceRegistration("remote", newRemoteDataService, types.LifetimeSingleton),
		registration.NamedServiceRegistration("local", newLocalDataService, types.LifetimeTransient))

	resolver := resolving.NewResolver(registry)
	scopedContext := resolving.NewScopedContext(t.Context())

	namedServiceFactory, _ := resolving.ResolveRequiredService[func(string) (dataService, error)](scopedContext, resolver)
	remote, _ := namedServiceFactory("remote")
	local, _ := namedServiceFactory("local")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, namedServiceFactory)
	assert.NotNil(t, remote)
	assert.NotNil(t, local)
}

func Test_Registry_register_named_service_consume_factory(t *testing.T) {

	// Arrange
	ctx := t.Context()

	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, newControllerWithNamedServiceFactory)
	_ = features.RegisterNamed[dataService](ctx, registry,
		registration.NamedServiceRegistration("remote", newRemoteDataService, types.LifetimeSingleton),
		registration.NamedServiceRegistration("local", newLocalDataService, types.LifetimeTransient))

	resolver := resolving.NewResolver(registry)
	scopedContext := resolving.NewScopedContext(ctx)

	// Act
	actual, err := resolving.ResolveRequiredService[*controllerWithNamedServices](scopedContext, resolver)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.NotNil(t, actual.remoteDataService)
	assert.NotNil(t, actual.localDataService)
}

func Test_Registry_register_named_service_resolve_all_named_services(t *testing.T) {

	// Arrange
	ctx := t.Context()

	registry := registration.NewServiceRegistry()
	_ = features.RegisterNamed[dataService](ctx, registry,
		registration.NamedServiceRegistration("remote", newRemoteDataService, types.LifetimeSingleton),
		registration.NamedServiceRegistration("local", newLocalDataService, types.LifetimeTransient))

	resolver := resolving.NewResolver(registry)
	scopedContext := resolving.NewScopedContext(ctx)

	// Act
	actual, err := resolving.ResolveRequiredServices[dataService](scopedContext, resolver)
	namedServiceFactory, _ := resolving.ResolveRequiredService[func(string) (dataService, error)](scopedContext, resolver)
	remote, _ := namedServiceFactory("remote")
	local, _ := namedServiceFactory("local")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 2, len(actual))
	assert.NotNil(t, remote)
	assert.Equal(t, "data from remote service", remote.FetchData())
	assert.NotNil(t, local)
	assert.Equal(t, "data from local service", local.FetchData())
}

func Test_Registry_register_named_service_resolve_all_named_services_as_list(t *testing.T) {

	// Arrange
	ctx := t.Context()

	registry := registration.NewServiceRegistry()
	_ = features.RegisterNamed[dataService](ctx, registry,
		registration.NamedServiceRegistration("remote", newRemoteDataService, types.LifetimeSingleton),
		registration.NamedServiceRegistration("local", newLocalDataService, types.LifetimeTransient))

	features.RegisterList[dataService](ctx, registry)

	resolver := resolving.NewResolver(registry)
	scopedContext := resolving.NewScopedContext(ctx)

	// Act
	actual, err := resolving.ResolveRequiredService[[]dataService](scopedContext, resolver)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, 2, len(actual))
}

type controllerWithNamedServices struct {
	remoteDataService dataService
	localDataService  dataService
}

func newControllerWithNamedServiceFactory(dataServiceFactory func(string) (dataService, error)) *controllerWithNamedServices {
	remote, _ := dataServiceFactory("remote")
	local, _ := dataServiceFactory("local")
	return &controllerWithNamedServices{
		remoteDataService: remote,
		localDataService:  local,
	}
}
