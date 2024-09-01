## What is Parsley?

Parsley is a powerful, easy-to-use reflection-based dependency injection package that integrates seamlessly into any Go application.

Just like the versatile green herb it's named after, Parsley enhances the flavor of your codebase by keeping your dependencies clean, organized, and easy to manage. Whether you're working on a small project or a complex system, Parsley brings the convenience of automated dependency injection to Go, inspired by the best practices from other languages. With features like automated lifetime management, proxy generation, and method interception, Parsley is the ingredient that makes your Go applications more maintainable and scalable.


## Why dependency injection for Golang?

While dependency injection (DI) is less common in Golang compared to other languages, the complexity it introduces can often be outweighed by the benefits it brings, especially in larger projects. As projects grow, the need for indirection and modularity becomes inevitable, and a DI framework like Parsley can seamlessly bridge the gap between dependency management and service activation. Parsley goes beyond resolving dependencies—it automates key aspects such as lifetime management, proxy generation, and interception, reducing boilerplate and enhancing maintainability. With these powerful features, why not let Parsley handle the heavy lifting for you?


## Features

- ✔️ Register types via constructor functions
- ✔️ Resolve objects by type (both interface and pointer type)
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
  - ✔️ Lazy loading objects by injecting dependencies as `Lazy[T]`
  - ✔️ Register factory functions to create instances of services based on input parameters provided at runtime
  - ⏳ Validate registered services; fail early during application startup if missing registrations are encountered
  - ✔️ Provide instances for non-registered types, use `ResolveWithOptions[T]` insted of `Resolve[T]`
- ✔️ Support multiple service registrations for the same interface
  - ✔️ Register named services (mutiple services), resolve via `func(key string) T`
  - ✔️ Resolve services as list (default)
- ✔️ Support for proxy types via code generation
  - ✔️ Proxies can be consumed as drop-in replacements for target services
  - ✔️ Proxies are extensible via `MethodInterceptor` services
- ⏳ Support sub-scopes
  - ⏳ Automatic clean-up


✔️ Already available | ❌ Not supported | ⏳ On schedule to be developed


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