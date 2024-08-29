package bootstrap

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type infrastructure struct {
	registry types.ServiceRegistry
	resolver types.Resolver
	app      Application
}

var parsley infrastructure

func RunParsleyApplication(cxt context.Context, appFactoryFunc any, configure ...types.ModuleFunc) error {

	registry := registration.NewServiceRegistry()
	registry.Register(appFactoryFunc, types.LifetimeSingleton)
	for _, m := range configure {
		m(registry)
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
