package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"sync"
)

type dependencyInfo struct {
	registration types.ServiceRegistration
	instance     interface{}
	children     []types.DependencyInfo
	consumer     types.DependencyInfo
	m            *sync.RWMutex
}

var _ types.DependencyInfo = &dependencyInfo{}

// NewDependencyInfo creates a new instance of types.DependencyInfo with the provided service registration, instance, and parent dependency.
func NewDependencyInfo(registration types.ServiceRegistration, instance interface{}, consumer types.DependencyInfo) types.DependencyInfo {
	return &dependencyInfo{
		registration: registration,
		instance:     instance,
		children:     make([]types.DependencyInfo, 0),
		consumer:     consumer,
		m:            &sync.RWMutex{},
	}
}

// Consumer returns the parent dependency for the current dependency info.
func (d *dependencyInfo) Consumer() types.DependencyInfo {
	return d.consumer
}

// Registration retrieves the service registration associated with this dependency info.
func (d *dependencyInfo) Registration() types.ServiceRegistration {
	return d.registration
}

// AddRequiredServiceInfo adds a child dependency info to the current dependency info to track service requirements.
func (d *dependencyInfo) AddRequiredServiceInfo(child types.DependencyInfo) {
	d.m.Lock()
	defer d.m.Unlock()
	d.children = append(d.children, child)
}

// RequiredServiceTypes fetches the types of services needed by the current dependency info from its registration.
func (d *dependencyInfo) RequiredServiceTypes() []types.ServiceType {
	return d.registration.RequiredServiceTypes()
}

// RequiredServices retrieves instances of services required by the current dependency.
func (d *dependencyInfo) RequiredServices() ([]interface{}, error) {
	d.m.RLock()
	defer d.m.RUnlock()
	instances := make([]interface{}, 0)
	for _, child := range d.children {
		instances = append(instances, child.Instance())
	}
	return instances, nil
}

// CreateInstance initializes and returns the service instance for this dependency.
// It resolves required services and uses the activator to create the instance.
func (d *dependencyInfo) CreateInstance() (interface{}, error) {
	if d.instance != nil {
		return d.instance, nil
	}
	resolvedDependencies, _ := d.RequiredServices()
	instance, err := d.registration.InvokeActivator(resolvedDependencies...)
	if err != nil {
		return nil, err
	}
	d.instance = instance
	return d.instance, nil
}

// SetInstance sets the instance for the current dependency if not already set, otherwise returns an error.
// This ensures that each dependency can only have one instance assigned, preventing unexpected behavior in the dependency lifecycle.
func (d *dependencyInfo) SetInstance(instance interface{}) error {
	d.m.Lock()
	defer d.m.Unlock()
	if d.instance != nil {
		return types.NewDependencyError(types.ErrorInstanceAlreadySet)
	}
	d.instance = instance
	return nil
}

// HasInstance checks if an instance has already been created for the dependency.
func (d *dependencyInfo) HasInstance() bool {
	return d.instance != nil
}

// Instance returns the existing service instance for this dependency info.
// Useful to retrieve a cached instance without (re-)initializing it, otherwise use CreateInstance instead.
func (d *dependencyInfo) Instance() interface{} {
	return d.instance
}

// ServiceTypeName returns the name of the service type associated with this dependency info.
// This helps in identifying the type of service being registered or retrieved.
func (d *dependencyInfo) ServiceTypeName() string {
	return d.registration.ServiceType().Name()
}
