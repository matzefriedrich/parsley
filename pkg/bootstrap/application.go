package bootstrap

import (
	"context"
	"errors"

	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type infrastructure struct {
	registry types.ServiceRegistry
	resolver types.Resolver
	app      Application
}

const (
	ErrorCannotRegisterAppFactory = "cannot register application factory"
)

// ErrCannotRegisterAppFactory is returned when the application factory cannot be registered, indicating an issue with the bootstrap process.
var (
	ErrCannotRegisterAppFactory = errors.New(ErrorCannotRegisterAppFactory)
)

// RunParsleyApplication initializes and runs the Parsley application lifecycle.
// It registers the application factory, configures additional modules, resolves the main application instance, and invokes its Run method.
func RunParsleyApplication(cxt context.Context, appFactoryFunc any, configure ...types.ModuleFunc) error {

	registry := registration.NewServiceRegistry()
	registerErr := registry.Register(appFactoryFunc, types.LifetimeSingleton)
	if registerErr != nil {
		bootstrapErr := &types.ParsleyError{Msg: ErrorCannotRegisterAppFactory}
		types.WithCause(registerErr)(bootstrapErr)
		return bootstrapErr
	}
	for _, m := range configure {
		_ = m(registry)
	}

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(cxt)
	app, appErr := resolving.ResolveRequiredService[Application](ctx, resolver)
	if appErr != nil {
		activationErr := &types.ParsleyError{Msg: "failed to activate application"}
		types.WithCause(appErr)(activationErr)
		return activationErr
	}

	parsley := infrastructure{
		registry: registry,
		resolver: resolver,
		app:      app,
	}

	appContext := context.WithValue(ctx, core.ContextKey("__parsley-infrastructure"), parsley)

	return app.Run(appContext)
}
