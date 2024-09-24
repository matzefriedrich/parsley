package bootstrap

import "context"

// Application provides an abstract interface for creating and running an application. It primarily facilitates the use of dependency injection for resolving services and the managing application lifecycle.
type Application interface {
	Run(context.Context) error
}
