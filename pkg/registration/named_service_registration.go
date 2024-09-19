package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
)

// NamedServiceRegistrationFunc defines a function that returns a service name, its activator function, and its lifetime scope. This type supports the internal infrastructure.
type NamedServiceRegistrationFunc func() (name string, activatorFunc any, scope types.LifetimeScope)

// NamedServiceRegistration registers a service with a specified name, activator function, and lifetime scope.
func NamedServiceRegistration(name string, activatorFunc any, scope types.LifetimeScope) NamedServiceRegistrationFunc {
	return func() (string, any, types.LifetimeScope) {
		return name, activatorFunc, scope
	}
}
