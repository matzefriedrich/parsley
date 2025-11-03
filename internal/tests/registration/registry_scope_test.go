package registration

import (
	"testing"

	"github.com/matzefriedrich/parsley/internal/tests/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_CreateScope_inherits_registered_types(t *testing.T) {

	// Arrange
	greeterFunc := func() features.Greeter {
		return features.NewGreeterMock()
	}

	globalRegistry := registration.NewServiceRegistry()
	_ = globalRegistry.Register(greeterFunc, types.LifetimeTransient)

	// Act
	scopedRegistry := globalRegistry.CreateScope()
	_ = scopedRegistry.Register(newTestService, types.LifetimeTransient)

	resolver := resolving.NewResolver(scopedRegistry)
	scopedContext := resolving.NewScopedContext(t.Context())
	actual, _ := resolving.ResolveRequiredService[*testService](scopedContext, resolver)

	// Assert
	assert.NotNil(t, scopedRegistry)
	assert.NotNil(t, actual)

	greeterServiceType := types.MakeServiceType[features.Greeter]()
	testServiceType := types.MakeServiceType[*testService]()

	assert.True(t, globalRegistry.IsRegistered(greeterServiceType))
	assert.False(t, globalRegistry.IsRegistered(testServiceType))
	assert.True(t, scopedRegistry.IsRegistered(testServiceType))
	assert.True(t, scopedRegistry.IsRegistered(greeterServiceType))
}
