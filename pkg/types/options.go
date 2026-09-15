package types

// DefaultLifecycleGroup is the implicit lifecycle group assigned to services registered without an explicit group.
const DefaultLifecycleGroup = "default"

// registrationOptions holds the configuration gathered from LifecycleOption values.
type registrationOptions struct {
	lifecycleGroup    string
	lifecycleGroupSet bool
	multipleGroups    bool
}

// LifecycleOption configures a service registration with lifecycle settings.
type LifecycleOption func(*registrationOptions)

// InLifecycleGroup assigns the service to the given lifecycle group. The group name must not be empty; applying
// InLifecycleGroup more than once with different group names results in a registry error.
func InLifecycleGroup(name string) LifecycleOption {
	return func(o *registrationOptions) {
		o.lifecycleGroup = name
		o.lifecycleGroupSet = true
	}
}

// ApplyLifecycleOptions applies the given options and returns the deduced lifecycle group. It returns a RegistryError
// when an empty group name was used or when multiple distinct groups were assigned. This function supports internal
// infrastructure.
func ApplyLifecycleOptions(options ...LifecycleOption) (string, error) {
	o := &registrationOptions{}
	for _, apply := range options {
		previousGroup := o.lifecycleGroup
		wasSet := o.lifecycleGroupSet
		apply(o)
		if wasSet && o.lifecycleGroupSet && previousGroup != o.lifecycleGroup {
			o.multipleGroups = true
		}
	}
	if o.lifecycleGroupSet && o.lifecycleGroup == "" {
		return "", NewRegistryError(ErrorEmptyLifecycleGroupName)
	}
	if o.multipleGroups {
		return "", NewRegistryError(ErrorMultipleLifecycleGroups)
	}
	return o.lifecycleGroup, nil
}