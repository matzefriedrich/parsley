[![CI](https://github.com/matzefriedrich/parsley/actions/workflows/go.yml/badge.svg)](https://github.com/matzefriedrich/parsley/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/matzefriedrich/parsley/badge.svg)](https://coveralls.io/github/matzefriedrich/parsley)
[![Go Reference](https://pkg.go.dev/badge/github.com/matzefriedrich/parsley.svg)](https://pkg.go.dev/github.com/matzefriedrich/parsley)
[![Go Report Card](https://goreportcard.com/badge/github.com/matzefriedrich/parsley)](https://goreportcard.com/report/github.com/matzefriedrich/parsley)
![License](https://img.shields.io/github/license/matzefriedrich/parsley)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matzefriedrich/parsley)
![GitHub Release](https://img.shields.io/github/v/release/matzefriedrich/parsley?include_prereleases)

## What is Parsley?

Parsley is a powerful, easy-to-use reflection-based dependency injection package that integrates seamlessly into any Go application.

Just like the versatile green herb it's named after, Parsley enhances the flavor of your codebase by keeping your dependencies clean, organized, and easy to manage. Whether you're working on a small project or a complex system, Parsley brings the convenience of automated dependency injection to Go, inspired by the best practices from other languages. With features like automated lifetime management, proxy generation, and method interception, Parsley is the ingredient that makes your Go applications more maintainable and scalable.


### Why dependency injection for Golang?

While dependency injection (DI) is less common in Golang compared to other languages, the complexity it introduces can often be outweighed by the benefits it brings, especially in larger projects. As projects grow, the need for indirection and modularity becomes inevitable, and a DI framework like Parsley can seamlessly bridge the gap between dependency management and service activation. Parsley goes beyond resolving dependencies—it automates key aspects such as lifetime management, proxy generation, and interception, reducing boilerplate and enhancing maintainability. With these powerful features, why not let Parsley handle the heavy lifting for you?


## Key Features

### Type Registration

- **Constructor Functions:** Register types via constructor functions.
- **Resolve by Type:** Resolve objects by both interface and pointer types.
  - **Constructor Injection:** Supported.
  - **Field Injection:** ⏳ Will be supported via struct tags (in progress).
  - **Setter Method Injection:** ❌ Not supported.
  - **Safe Casts:** Convenience function `ResolveRequiredService[T]` to resolve and safely cast objects.

- **Lifetime Management:** Register types with different lifetimes.
  - **Singleton:** Available.
  - **Register Instances as Singletons:** Use `RegisterInstance[T]` (where `T` must be an interface type).
  - **Scoped:** Requires a `NewScopedContext(context.Background)`. Use `RegisterScoped()`.
  - **Transient:** Available for objects with a transient lifetime.

### Module & On-Demand Registrations

- **Modular Registrations:** Bundle type registrations as modules, register them via `RegisterModule` as a unit.

- **On-Demand Resolution:** Resolve objects only when needed.
  - **Custom Factories:** Allow consumption of `Resolver` instead of custom factories.
  - **Lazy Loading:** Inject dependencies lazily using `Lazy[T]`.
  - **Factory Functions:** Register factories to create service instances dynamically at runtime.
  - **Service Validation:** ⏳ Planned feature to validate services during startup and fail early if missing registrations are found (in progress).
  - **Non-Registered Types:** Resolve non-registered types using `ResolveWithOptions[T]`.
  - **Override Type Registrations:** Provide custom service instances when resolving service instances.

### Multiple Registrations for the Same Interface

- **Named Services:** Register and resolve multiple services for the same interface, via `func(key string) T`.
- **Service Lists:** Resolve services as a list. Enable list dependencies via `RegisterList`.

### Proxy Type Support

- **Proxies as Drop-in Replacements:** Generate proxy types that can be used in place of target services.
- **Extensible Proxies:** Proxies are extensible with `MethodInterceptor` services.

### Configurable Mocks

- **Generate Mocks:** Generate configurable mocks via `//go:parsley-cli generate mocks` to boost automated testing.


## Usage

Getting started with Parsley is as simple as adding it to your project and letting it take care of your dependencies. Here’s how you can use Parsley to wire up a small example application.

### Install Parsley

First, add Parsley to your project:

````sh
go get github.com/matzefriedrich/parsley
````

For more advanced features, you can install the `parsley-cli` utility:

````sh
go install github.com/matzefriedrich/parsley/cmd/parsley-cli
````

### Creating a Parsley-powered "Hello-World" application

Imagine you're building a service that greets users and logs those greetings. Instead of manually wiring everything together, Parsley can handle the setup. Start by defining the interfaces for your services:

````go
// Greeter defines a service that greets users.
type Greeter interface {
    SayHello(name string) string
}

// Logger defines a service that logs messages.
type Logger interface {
    Log(message string)
}
````

Now, create the implementations for these interfaces:

````go
type greeterServiceImpl struct {
    logger LoggerService
}

func (g *greeterServiceImpl) SayHello(name string) string {
    greeting := "Hello, " + name + "!"
    g.logger.Log("Generated greeting: " + greeting)
    return greeting
}

type loggerServiceImpl struct{}

func (l *loggerServiceImpl) Log(message string) {
    fmt.Println("Log:", message)
}
````

To wire everything together, define constructor functions:

````go
func NewGreetService(logger LoggerService) GreetService {
    return &greetServiceImpl{logger: logger}
}

func NewLoggerService() LoggerService {
    return &loggerServiceImpl{}
}
````

With Parsley, registering these services and resolving them is straightforward:

````go
func main() {
    registry := registration.NewServiceRegistry()

    // Register services with their lifetimes
    registration.RegisterSingleton(registry, NewLoggerService)
    registration.RegisterScoped(registry, NewGreetService)

    // Create a resolver
    resolver := resolving.NewResolver(registry)

    // Resolve and use the Greeter service instance
    scope := resolving.NewScopedContext(context.Background())
    greeter, _ := resolving.ResolveRequiredService[Greeter](resolver, scope)

    // Use the service
    fmt.Println(greeter.SayHello("World"))
}
````

When you run this application, you’ll see that Parsley automatically handles the dependencies for you:

````sh
Log: Generated greeting: Hello, World!
Hello, World!
````

### What Just Happened?

In this example, you defined two services, `Greeter` and `Logger`. You then registered these services with Parsley, specifying that `Logger` should have a singleton lifetime while `Greeter` should be scoped. Parsley injected a `Logger` instance into the `Greeter` instance during creation, ensuring everything was wired up correctly.

By using Parsley, you avoid the hassle of manual dependency wiring, keeping your code clean and focused on business logic. But that is not all - there are more features to explore: head over to the official docs at https://matzefriedrich.github.io/parsley-docs/ to find more usage examples.


---

Copyright 2024 - Matthias Friedrich
