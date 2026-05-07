package resolving

import (
	"context"
	"errors"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Resolver_verify_different_activator_function_signatures(t *testing.T) {

	t.Run("Support (T, error) success", func(t *testing.T) {
		// Arrange
		registry := registration.NewServiceRegistry()
		activatorFuncWithError := func() (ServiceA, error) {
			return &serviceA{val: "A"}, nil
		}
		_ = registry.Register(activatorFuncWithError, types.LifetimeTransient)
		resolver := resolving.NewResolver(registry)

		ctx := t.Context()

		// Act
		actual, err := resolving.ResolveRequiredService[ServiceA](ctx, resolver)

		assert.NoError(t, err)
		assert.Equal(t, "A", actual.DoA())
	})

	t.Run("Support (T, error) failure", func(t *testing.T) {
		// Arrange
		registry := registration.NewServiceRegistry()
		failingActivatorFuncWithError := func() (ServiceA, error) {
			return nil, errors.New("failed to create A")
		}
		_ = registry.Register(failingActivatorFuncWithError, types.LifetimeTransient)
		resolver := resolving.NewResolver(registry)

		ctx := t.Context()

		// Act
		_, err := resolving.ResolveRequiredService[ServiceA](ctx, resolver)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, errors.Unwrap(err).Error(), "failed to create A")
	})

	t.Run("Support (context.Context, ...) T", func(t *testing.T) {
		// Arrange
		registry := registration.NewServiceRegistry()
		activatorFuncWithContext := func(activatorCtx context.Context) ServiceA {
			val := activatorCtx.Value("test-key").(string)
			return &serviceA{val: val}
		}
		_ = registry.Register(activatorFuncWithContext, types.LifetimeTransient)
		resolver := resolving.NewResolver(registry)

		ctx := context.WithValue(t.Context(), "test-key", "context-A")

		// Act
		actual, err := resolving.ResolveRequiredService[ServiceA](ctx, resolver)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "context-A", actual.DoA())
	})

	t.Run("Support (context.Context, ...) (T, error)", func(t *testing.T) {
		// Arrange
		registry := registration.NewServiceRegistry()
		failingActivatorFuncWithContext := func(ctx context.Context) (ServiceA, error) {
			val := ctx.Value("test-key").(string)
			return &serviceA{val: val}, nil
		}
		_ = registry.Register(failingActivatorFuncWithContext, types.LifetimeTransient)
		resolver := resolving.NewResolver(registry)

		ctx := context.WithValue(t.Context(), "test-key", "context-A-error")

		// Act
		actual, err := resolving.ResolveRequiredService[ServiceA](ctx, resolver)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "context-A-error", actual.DoA())
	})

	t.Run("Support standard T", func(t *testing.T) {
		// Arrange
		registry := registration.NewServiceRegistry()
		activatorFunc := func() ServiceA {
			return &serviceA{val: "A"}
		}
		_ = registry.Register(activatorFunc, types.LifetimeTransient)
		resolver := resolving.NewResolver(registry)

		ctx := context.Background()

		// Act
		svc, err := resolving.ResolveRequiredService[ServiceA](ctx, resolver)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "A", svc.DoA())
	})
}

type ServiceA interface {
	DoA() string
}

type serviceA struct {
	val string
}

func (s *serviceA) DoA() string {
	return s.val
}
