package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
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

func NewDependencyInfo(registration types.ServiceRegistration, instance interface{}, consumer types.DependencyInfo) types.DependencyInfo {
	return &dependencyInfo{
		registration: registration,
		instance:     instance,
		children:     make([]types.DependencyInfo, 0),
		consumer:     consumer,
		m:            &sync.RWMutex{},
	}
}

func (d *dependencyInfo) Consumer() types.DependencyInfo {
	return d.consumer
}

func (d *dependencyInfo) Registration() types.ServiceRegistration {
	return d.registration
}

func (d *dependencyInfo) AddRequiredServiceInfo(child types.DependencyInfo) {
	d.m.Lock()
	defer d.m.Unlock()
	d.children = append(d.children, child)
}

func (d *dependencyInfo) RequiredServiceTypes() []reflect.Type {
	return d.registration.RequiredServiceTypes()
}

func (d *dependencyInfo) RequiredServices() ([]interface{}, error) {
	d.m.RLock()
	defer d.m.RUnlock()
	instances := make([]interface{}, 0)
	for _, child := range d.children {
		instances = append(instances, child.Instance())
	}
	return instances, nil
}

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

func (d *dependencyInfo) SetInstance(instance interface{}) error {
	d.m.Lock()
	defer d.m.Unlock()
	if d.instance != nil {
		return types.NewDependencyError(types.ErrorInstanceAlreadySet)
	}
	d.instance = instance
	return nil
}

func (d *dependencyInfo) HasInstance() bool {
	return d.instance != nil
}

func (d *dependencyInfo) Instance() interface{} {
	return d.instance
}

func (d *dependencyInfo) ServiceTypeName() string {
	return d.registration.ServiceType().Name()
}
