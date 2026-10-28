package bootstrap

import (
	"context"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

type singletonService struct {
	disposed bool
}

func (s *singletonService) Dispose(ctx context.Context) error {
	s.disposed = true
	return nil
}

type scopedService struct {
	disposed bool
}

func (s *scopedService) Dispose(ctx context.Context) error {
	s.disposed = true
	return nil
}

func Test_RunParsleyApplication_disposes_singleton_and_scoped_services(t *testing.T) {

	// Arrange
	var singleton *singletonService
	var scoped *scopedService

	appFactory := func(s *singletonService, resolver types.Resolver) bootstrap.Application {
		singleton = s
		return &testApp{
			RunFunc: func(ctx context.Context) error {
				var err error
				scoped, err = resolving.ResolveRequiredService[*scopedService](ctx, resolver)
				return err
			},
		}
	}

	// Act
	err := bootstrap.RunParsleyApplication(t.Context(), appFactory, func(registry types.ServiceRegistry) error {
		_ = registration.RegisterSingleton(registry, func() *singletonService { return &singletonService{} })
		_ = registration.RegisterScoped(registry, func() *scopedService { return &scopedService{} })
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, singleton)
	assert.NotNil(t, scoped)
	assert.True(t, singleton.disposed, "singleton should be disposed")
	assert.True(t, scoped.disposed, "scoped should be disposed")
}

type orderedScopedService struct {
	id      string
	history *[]string
}

func (s *orderedScopedService) Dispose(_ context.Context) error {
	*s.history = append(*s.history, s.id)
	return nil
}

func Test_RunParsleyApplication_teardown_honours_lifecycle_groups(t *testing.T) {

	// Arrange
	history := make([]string, 0)

	appFactory := func(resolver types.Resolver) bootstrap.Application {
		return &testApp{
			RunFunc: func(ctx context.Context) error {
				_, _ = resolving.ResolveRequiredServices[*orderedScopedService](ctx, resolver)
				return nil
			},
		}
	}

	// Act
	err := bootstrap.RunParsleyApplication(t.Context(), appFactory, func(registry types.ServiceRegistry) error {
		registry.SetTeardownOrder("transport", "infrastructure")
		_ = registration.RegisterScopedWithOptions(registry, func() *orderedScopedService {
			return &orderedScopedService{id: "http", history: &history}
		}, types.InLifecycleGroup("transport"))
		_ = registration.RegisterScopedWithOptions(registry, func() *orderedScopedService {
			return &orderedScopedService{id: "db", history: &history}
		}, types.InLifecycleGroup("infrastructure"))
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"http", "db"}, history)
}
