package registration

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_SetTeardownOrder_normalizes_and_deduplicates(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	registry.SetTeardownOrder("transport", "infrastructure", "", "infrastructure", "db")

	// Assert
	assert.Equal(t, []string{"transport", "infrastructure", "db"}, registry.GetTeardownOrder())
}

func Test_Registry_SetTeardownOrder_called_after_registration_is_allowed(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, newGroupableService)

	// Act
	registry.SetTeardownOrder("transport")

	// Assert
	assert.Equal(t, []string{"transport"}, registry.GetTeardownOrder())
}

func Test_Registry_SetTeardownOrder_last_call_replaces(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	registry.SetTeardownOrder("a", "b")
	registry.SetTeardownOrder("c")

	// Assert
	assert.Equal(t, []string{"c"}, registry.GetTeardownOrder())
}

func Test_Registry_SetTeardownOrder_without_arguments_resets(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("a", "b")

	// Act
	registry.SetTeardownOrder()

	// Assert
	assert.Empty(t, registry.GetTeardownOrder())
}

func Test_Registry_GetTeardownOrder_returns_a_copy(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("a", "b")

	// Act
	order := registry.GetTeardownOrder()
	order[0] = "changed"

	// Assert
	assert.Equal(t, []string{"a", "b"}, registry.GetTeardownOrder())
}

func Test_Registry_CreateScope_inherits_the_teardown_order(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("transport", "infrastructure")

	// Act
	scoped := registry.CreateScope()

	// Assert
	assert.Equal(t, []string{"transport", "infrastructure"}, scoped.GetTeardownOrder())
}

func Test_Registry_CreateLinkedRegistry_inherits_the_teardown_order(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("transport")

	// Act
	linked := registry.CreateLinkedRegistry()

	// Assert
	assert.Equal(t, []string{"transport"}, linked.GetTeardownOrder())
}
