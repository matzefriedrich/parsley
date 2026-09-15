package resolving

import (
	"context"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

type disposableService struct {
	disposed bool
}

func (s *disposableService) Dispose(ctx context.Context) error {
	s.disposed = true
	return nil
}

func newDisposableService() *disposableService {
	return &disposableService{}
}

func Test_Resolver_Shutdown_disposes_singleton_services(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, newDisposableService)

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	service, _ := resolving.ResolveRequiredService[*disposableService](ctx, r)

	// Act
	err := r.Shutdown(ctx)

	// Assert
	assert.NoError(t, err)
	assert.True(t, service.disposed)
}

func Test_DisposeScope_disposes_scoped_services(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterScoped(registry, newDisposableService)

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	service, _ := resolving.ResolveRequiredService[*disposableService](ctx, r)

	// Act
	err := resolving.DisposeScope(ctx)

	// Assert
	assert.NoError(t, err)
	assert.True(t, service.disposed)
}

func Test_Resolver_Shutdown_does_not_dispose_transient_services(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterTransient(registry, newDisposableService)

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	service, _ := resolving.ResolveRequiredService[*disposableService](ctx, r)

	// Act
	err := r.Shutdown(ctx)

	// Assert
	assert.NoError(t, err)
	assert.False(t, service.disposed)
}

type orderedDisposableService struct {
	id      int
	history *[]int
}

func (s *orderedDisposableService) Dispose(_ context.Context) error {
	*s.history = append(*s.history, s.id)
	return nil
}

func Test_Resolver_Shutdown_disposes_in_reverse_order(t *testing.T) {

	// Arrange
	history := make([]int, 0)
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, func() *orderedDisposableService {
		return &orderedDisposableService{id: 1, history: &history}
	})
	_ = registration.RegisterSingleton(registry, func() *orderedDisposableService {
		return &orderedDisposableService{id: 2, history: &history}
	})

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*orderedDisposableService](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []int{2, 1}, history)
}

type groupOrderedDisposable struct {
	id      string
	history *[]string
}

func (s *groupOrderedDisposable) Dispose(_ context.Context) error {
	*s.history = append(*s.history, s.id)
	return nil
}

func registerGroupedDisposablePair(registry types.ServiceRegistry, history *[]string) {
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "http", history: history}
	}, types.InLifecycleGroup("transport"))
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "db", history: history}
	}, types.InLifecycleGroup("infrastructure"))
}

func Test_Resolver_Shutdown_disposes_groups_in_declared_order(t *testing.T) {

	// Arrange
	history := make([]string, 0)
	registry := registration.NewServiceRegistry()
	registerGroupedDisposablePair(registry, &history)
	registry.SetTeardownOrder("transport", "infrastructure")

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*groupOrderedDisposable](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []string{"http", "db"}, history)
}

func Test_Resolver_Shutdown_disposes_within_group_in_reverse_resolution_order(t *testing.T) {

	// Arrange
	history := make([]string, 0)
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("db")
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "conn1", history: &history}
	}, types.InLifecycleGroup("db"))
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "conn2", history: &history}
	}, types.InLifecycleGroup("db"))

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*groupOrderedDisposable](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []string{"conn2", "conn1"}, history)
}

func Test_Resolver_Shutdown_disposes_unlisted_groups_last(t *testing.T) {

	// Arrange
	history := make([]string, 0)
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("transport")
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "http", history: &history}
	}, types.InLifecycleGroup("transport"))
	_ = registration.RegisterSingleton(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "metrics", history: &history}
	})

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*groupOrderedDisposable](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []string{"http", "metrics"}, history)
}

func Test_Resolver_Shutdown_explicit_default_group_in_declared_order(t *testing.T) {

	// Arrange
	history := make([]string, 0)
	registry := registration.NewServiceRegistry()
	registry.SetTeardownOrder("default", "transport")
	_ = registration.RegisterSingleton(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "metrics", history: &history}
	})
	_ = registration.RegisterSingletonWithOptions(registry, func() *groupOrderedDisposable {
		return &groupOrderedDisposable{id: "http", history: &history}
	}, types.InLifecycleGroup("transport"))

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*groupOrderedDisposable](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []string{"metrics", "http"}, history)
}

func Test_Resolver_Shutdown_skips_empty_listed_group(t *testing.T) {

	// Arrange
	history := make([]string, 0)
	registry := registration.NewServiceRegistry()
	registerGroupedDisposablePair(registry, &history)
	registry.SetTeardownOrder("unused", "transport", "infrastructure")

	r := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	_, _ = resolving.ResolveRequiredServices[*groupOrderedDisposable](ctx, r)
	_ = r.Shutdown(ctx)

	// Assert
	assert.Equal(t, []string{"http", "db"}, history)
}
