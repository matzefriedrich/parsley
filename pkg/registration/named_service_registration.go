package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
)

type NamedServiceRegistrationFunc func() (name string, activatorFunc any, scope types.LifetimeScope)

func NamedServiceRegistration(name string, activatorFunc any, scope types.LifetimeScope) NamedServiceRegistrationFunc {
	return func() (string, any, types.LifetimeScope) {
		return name, activatorFunc, scope
	}
}
