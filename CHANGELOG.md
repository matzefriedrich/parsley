# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - 



## v0.3.0 - 2024-07-12

### Added

* Service registrations can be bundled in a `ModuleFunc` to register related types as a unit
* Service registry accepts object instances as singleton service registrations
* Adds the `ResolveRequiredService[T]` convenience function that resolves and safe-casts objects
* Registers resolver instance with the registry so that the `Resolver` object can be injected into factory and constructor methods
* The resolver can now accept instances of non-registered types via the `ResolveWithOptions[T]` method
* `ServiceRegistry` has new methods for creating linked and scoped registry objects (which share the same `ServiceIdSequence`). Scoped registries inherit all parent service registrations, while linked registries are empty. See `CreateLinkedRegistry` and `CreateScope` methods.
  
### Changed

* A `ServiceRegistryAccessor` is no longer a `ServiceRegisty`, it is the other way around
* The creation of service registrations and type activators has been refactored; see `activator.go` and `service_registration.go` modules
* Multiple registries can be grouped with `NewMultiRegistryAccessor` to simplify the lookup of service registrations from linked registries. The resolver uses this accessor type to merge registered service types with object instances for unregistered types.


## v0.2.0 - 2024-07-11

### Added

* Resolver can now detect circular dependencies
* Adds helpers to register services with a certain lifetime scope

### Changed

* Registry rejects non-interface types

### Fixes

* Fixes error wrapping in custom error types
* Improves error handling for service registry and resolver


## v0.1.0 - 2024-07-10

### Added

* Adds service registry; the registry can map interfaces to implementation types via constructor functions
* Assign lifetime behaviour to services (singleton, scoped, or transient)
* Adds resolver (container) service
