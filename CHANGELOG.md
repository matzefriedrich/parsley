# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.7.0] - 2024-08-05

* Adds the `RegisterLazy[T]` method to register lazy service factories. Use the type `Lazy[T]` to consume a lazy service dependency and call the `Value() T` method on the lazy factory to request the actual service instance. The factory will create the service instance upon the first request, cache it, and return it for subsequent calls using the `Value` method.

## [v0.6.1] - 2024-07-30

### Changes

* Registers named services as transient services to resolve them also as a list of services (like services without a name). Changes the `createResolverRegistryAccessor` method so temporary registrations are selected first (and shadow permanent registrations). This behavior can also be leveraged in `ResolverOptionsFunc` to shadow other registrations when resolving instances via `ResolveWithOptions.`


## [v0.6.0] - 2024-07-26

### Added 

* Adds the `Activate[T]` method which can resolve an instance from an unregistered activator func.
* Allows registration and activation of pointer types (to not enforce usage of interfaces as abstractions).
* Adds the `RegisterNamed[T]` method to register services of the same interface and allow to resolve them by name.

### Changes

* Renames the `ServiceType[T]` method to `MakeServiceType[T]`; a service type represents now the reflected type and its name (which makes debugging and understanding service dependencies much easier).
* Replaces all usages of `reflect.Type` by `ServiceType` in all Parsley interfaces.
* Changes the `IsSame` method of the `ServiceRegistration` type; service registrations of type function are always treated as different service types.

### Fixes

* Fixes a bug in the `detectCircularDependency` function which could make the method get stuck in an infinite loop.


## v0.5.0 - 2024-07-16

### Added

* The service registry now accepts multiple registrations for the same interface (changes internal data structures to keep track of registrations; see `ServiceRegistrationList`).
* Adds the `ResolveRequiredServices[T]` convenience function to resolve all service instances; `ResolveRequiredService[T]` can resolve a single service but will return an error if service registrations are ambiguous.

### Changed

* Extends the resolver to handle multiple service registrations per interface type. The resolver returns resolved objects as a list. 


## v0.4.0 - 2024-07-13

### Added

* Support for factory functions to create instances of services based on input parameters provided at runtime

### Changed

* Reorganizes the whole package structure; adds sub-packages for `registration` and `resolving`. A bunch of types that support the inner functionality of the package have been moved to `internal.`.
* Integration tests are moved to the `internal` package.


## v0.3.0 - 2024-07-12

### Added

* Service registrations can be bundled in a `ModuleFunc` to register related types as a unit.
* The service registry accepts object instances as singleton service registrations.
* Adds the `ResolveRequiredService[T]` convenience function that resolves and safe-casts objects.
* Registers resolver instance with the registry so that the `Resolver` object can be injected into factory and constructor methods.
* The resolver can now accept instances of non-registered types via the `ResolveWithOptions[T]` method.
* `ServiceRegistry` has new methods for creating linked and scoped registry objects (which share the same `ServiceIdSequence`). Scoped registries inherit all parent service registrations, while linked registries are empty. See `CreateLinkedRegistry` and `CreateScope` methods.
  
### Changed

* A `ServiceRegistryAccessor` is no longer a `ServiceRegisty`, it is the other way around.
* The creation of service registrations and type activators has been refactored; see `activator.go` and `service_registration.go` modules.
* Multiple registries can be grouped with `NewMultiRegistryAccessor` to simplify the lookup of service registrations from linked registries. The resolver uses this accessor type to merge registered service types with object instances for unregistered types.


## v0.2.0 - 2024-07-11

### Added

* The resolver can now detect circular dependencies.
* Adds helpers to register services with a certain lifetime scope.

### Changed

* The registry rejects non-interface types.

### Fixes

* Fixes error wrapping in custom error types.
* Improves error handling for service registry and resolver.


## v0.1.0 - 2024-07-10

### Added

* Adds a service registry; the registry can map interfaces to implementation types via constructor functions.
* Assign lifetime behaviour to services (singleton, scoped, or transient).
* Adds a basic resolver (container) service.
