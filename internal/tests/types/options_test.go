package types

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_ApplyLifecycleOptions_no_options_returns_no_group(t *testing.T) {
	// Act
	group, err := types.ApplyLifecycleOptions()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, group)
}

func Test_ApplyLifecycleOptions_sets_the_group(t *testing.T) {
	// Act
	group, err := types.ApplyLifecycleOptions(types.InLifecycleGroup("infrastructure"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "infrastructure", group)
}

func Test_ApplyLifecycleOptions_empty_group_name_returns_error(t *testing.T) {
	// Act
	_, err := types.ApplyLifecycleOptions(types.InLifecycleGroup(""))

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func Test_ApplyLifecycleOptions_multiple_different_groups_returns_error(t *testing.T) {
	// Act
	_, err := types.ApplyLifecycleOptions(types.InLifecycleGroup("a"), types.InLifecycleGroup("b"))

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple lifecycle groups")
}

func Test_ApplyLifecycleOptions_same_group_applied_twice_last_wins(t *testing.T) {
	// Act
	group, err := types.ApplyLifecycleOptions(types.InLifecycleGroup("a"), types.InLifecycleGroup("a"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "a", group)
}
