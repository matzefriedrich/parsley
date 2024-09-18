package bootstrap

import (
	"context"
	"errors"
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

var (
	ErrCannotRegisterAppFactory = errors.New(ErrorCannotRegisterAppFactory)
)

var parsley infrastructure

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
	app, _ := resolving.ResolveRequiredService[Application](resolver, ctx)

	parsley = infrastructure{
		registry: registry,
		resolver: resolver,
		app:      app,
	}

	return app.Run(ctx)
}
