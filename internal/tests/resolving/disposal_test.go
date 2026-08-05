package resolving

import (
	"context"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
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
