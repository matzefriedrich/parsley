## What is Parsley?

Parsley is an easy-to-use reflection-based dependency injection package that fits into any Go application.

This dependency injection package may become your favorite ingredient for your Go applications. It is like the nifty green herb that fits well in various dishes across different cuisines. It not only adds to the taste, but it also charms the eye. In terms of wiring dependencies, it helps to keep things clean and organized. The parsley library is inspired by other dependency injection libraries I have used, which I always miss when working on Go projects.


## Why dependency injection for Golang?

Though dependency injection is less prevalent in Golang (compared to other languages), sometimes the burden of increased complexity by adopting a DI framework may be n-times less than the complexity of building larger projects without such a framework. The concept of indirection is usually given anyway in larger projects, so what is left to be done by a dependency injection framework is bridging the gap between dependency configuration and service activation. What you gain on top is automated lifetime behavior, to name one. So why not automate it?


## Features

- ✔️ Register types via constructor functions
- ✔️ Resolve objects by interface
  - ✔️ Constructor injection
  - ⏳ Injection via field initialization (requires annotation)
  - ❌ Injection via setter methods
  - ✔️ Convenience function to resolve and safe-cast objects: `ResolveRequiredService[T]`
- ✔️ Register types with a certain lifetime
  - ✔️ Singleton
  - ✔️ Register objects as singletons; use `RegisterInstance[T]` whereby `T` must be an interface type
  - ✔️ Scoped (requires a certain context `NewScopedContext(context.Background))`; use `RegisterScoped`)
  - ✔️ Transient
- ✔️ Bundle type registrations as modules to register them via `RegisterModule` as a unit
- ✔️ Resolve objects on-demand
  - ✔️ Allow consumption of `Resolver` in favor of custom factories
  - ⏳ Lazy loading objects by injecting dependencies as `Lazy[T]`
  - ✔️ Register factory functions to create instances of services based on input parameters provided at runtime
  - ⏳ Validate registered services; fail early during application startup if missing registrations are encountered
  - ✔️ Provide instances for non-registered types, use `ResolveWithOptions[T]` insted of `Resolve[T]`
- ✔️ Support multiple service registrations for the same interface
  - ⏳ Register named services (mutiple services), resolve via `func(key string) any`
  - ✔️ Resolve services as list (default)
- ⏳ Support sub-scopes
  - ⏳ Automatic clean-up


✔️ Already available | ❌ Not supported | ⏳ On schedule to be developed


## Usage

````sh
$ go get github.com/matzefriedrich/parsley
````


## Dependency mapping configuration

The dependency mapping configuration requires types, interfaces, and constructor methods—basically, the same things you need to wire dependencies manually. 

In parsley, constructor methods are the centerpiece that defines the mappings between abstractions (interfaces) and implementation types; required services are expressed as function parameters. The return type of a constructor method specifies the abstraction, whereby the method itself is responsible for creating the actual object instance, thus acting as the glue between implementation- and interface types.

### Learn Parsley by example

Parsley uses reflection to build a service registry and resolve objects at runtime. Use the `NewServiceRegistry` method to create a registry service that tracks service mapping and lifetime configuration. Use the `Register` method of the `ServiceRegistry` to register types via constructor methods.

Since the `ServiceRegistry` is quite generic, it may feel more natural to use the registration helper methods `RegisterSingleton`, `RegisterScoped`, `RegisterTransient`, or `RegisterInstance` to register types and configure a lifetime behavior.

````golang
registry := registration.NewServiceRegistry()
_ = registration.RegisterSingleton(registry, NewFoo)
_ = registration.RegisterScoped(registry, NewFooConsumer)
````

Once all service types are registered, a resolver service (the actual container) must be created. The resolver creates object instances and manages their lifetime following the service configuration. Use the `NewResolver` method and pass the registry instance that holds the type registrations. 

````golang
resolver := resolving.NewResolver(registry)
scope := internal.NewScopedContext(context.Background())
consumerInstance, _ := resolving.ResolveRequiredService[FooConsumer](resolver, scope)
````

The types and methods involved in the example above, are defined as follows:

````golang
type Foo interface {
    Bar()
}

type FooConsumer interface {
    FooBar()
}
````

Next, implementations for those services are defined as follows:

````golang
type foo struct{}

func (f *foo) Bar() {}

type fooConsumer struct {
    foo Foo
}

func (c *fooConsumer) FooBar() {
    c.foo.Bar()
}
````

The `Foo` service does not have any dependencies. Thus, the `NewFoo` constructor function has no parameters. It returns a new instance of `foo`; the resolver is not interested in the actual implementation type of `Foo`.

In contrast, the constructor function of `FooConsumer` requires a `Foo` object. Parsley builds a dependency graph for requested services, resolves those services first (respecting configured lifetimes), and then calls constructor methods with all required parameters.

````golang
func NewFoo() Foo {
    return &foo{}
}

func NewFooConsumer(foo Foo) FooConsumer {
    return &fooConsumer{
        foo: foo,
    }
}
````
