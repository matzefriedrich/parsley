# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.12] - 2025-05-02

* Bumps `golang.org/x/mod` from 0.23.0 to 0.24.0 [#47](https://github.com/matzefriedrich/parsley/pull/47)
* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.4.4 to 0.4.5 [#48](https://github.com/matzefriedrich/parsley/pull/48)


## [v1.0.11] - 2025-03-06

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.4.3 to 0.4.4 [#46](https://github.com/matzefriedrich/parsley/pull/46)


## [v1.0.10] - 2025-02-19

* Bumps `github.com/spf13/cobra from` 1.8.1 to 1.9.1 [#45](https://github.com/matzefriedrich/parsley/pull/45)


## [v1.0.9] - 2025-02-12

### Changed

* Bumps `golang.org/x/mod` from 0.22.0 to 0.23.0 [#44](https://github.com/matzefriedrich/parsley/pull/44)


## [v1.0.8] - 2025-01-23

### Changed

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.4.2 to 0.4.3 [#43](https://github.com/matzefriedrich/parsley/pull/44)


## [v1.0.7] - 2024-12-17

### Changed

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.4.1 to 0.4.2 [#42](https://github.com/matzefriedrich/parsley/pull/42)


## [v1.0.6] - 2024-12-11

### Changed

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.4.0 to 0.4.1 [#41](https://github.com/matzefriedrich/parsley/pull/41)


## [v1.0.5] - 2024-11-27

### Changed

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.3.2 to 0.4.0 [#40](https://github.com/matzefriedrich/parsley/pull/40)
* Alters command metadata (adds long descriptions)

## [v1.0.4] - 2024-11-18

### Changed

* Bumps `github.com/stretchr/testify` from 1.9.0 to 1.10.0 [#39](https://github.com/matzefriedrich/parsley/pull/39)
* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.3.1 to 0.3.2 [#38](https://github.com/matzefriedrich/parsley/pull/39)
* Bumps `golang.org/x/mod` from 0.21.0 to 0.22.0 [#37](https://github.com/matzefriedrich/parsley/pull/37)


## [v1.0.3] - 2024-10-08

### Changed

* Removed `hashicorp/go-version` dependency; added simple comparison functions to the `version.go` module instead. [#36](https://github.com/matzefriedrich/parsley/pull/36)


## [v1.0.2] - 2024-09-26

### Changed

* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.3.0 to 0.3.1 [#35](https://github.com/matzefriedrich/parsley/pull/35)
* Bumps `github.com/matzefriedrich/cobra-extensions` from 0.2.6 to 0.3.0, and adjusts import paths [#34](https://github.com/matzefriedrich/parsley/pull/34)


## [v1.0.1] - 2024-09-23

### Added

* Added new tests specifically for the service registration convenience functions. [#31](https://github.com/matzefriedrich/parsley/pull/31) 
* Added new tests for template functions and code generation. [#32](https://github.com/matzefriedrich/parsley/pull/32)
* Adds documentation texts.

### Changed

* Adds the `SupportsRegisterActivatorFunc`, which is used in the registration functions instead of `ServiceRegistry`. 

* The convenience functions for registering services with different lifetimes are moved to the `register_functions.go` module, for better organization and separation of concerns. [#31](https://github.com/matzefriedrich/parsley/pull/31)

* Multiple activator functions can be passed to the `RegisterTransient`, `RegisterSingleton`, and `RegisterScoped` convenience functions, allowing several services to be registered with a single method call. [#31](https://github.com/matzefriedrich/parsley/pull/31) 

### Fixed

* Several improvements and fixes related to handling interface and ellipsis parameters in reflection and code generation. [#32](https://github.com/matzefriedrich/parsley/pull/32)


## [v1.0.0] - 2024-09-21

### Added

* Added the `validator.go` module to the `registration` package, introducing a `Validator` service to verify service registrations [#30](https://github.com/matzefriedrich/parsley/pull/30):
  * Detects missing service dependencies.
  * Identifies circular service registrations, where services depend on themselves.


## [v0.10.1] - 2024-09-20

### Changed

* Added documentation comments to all exported types and methods.


## [v0.10.0] - 2024-09-19

### Changed

Increased overall test coverage for the Parsley library and CLI utility from 55% to 80%, improving reliability and confidence in the codebase.
  
### Added

* Introduced new tests for almost all packages and modules
  
### Fixed

* Fixed issues discovered during testing, particularly in `ParsleyAggregateError` and the `reflection` and `generator` packages.

### Changed

* Refactored the (internal) `reflection` and `generator` packages to enhance testability. Made minor internal adjustments to ensure better separation of concerns and improved testability. The updated type-model builder now supports pointer—and array-type parameters and result fields (the previous implementation could only handle scalar and simple array types).

* Introduced the `ParameterType` struct to store reflected parameter and field type information.


## [v0.9.3] - 2024-09-13

### Changed

* The `type_model_builder.go` module has been refactored to improve flexibility and maintain separation of concerns. By decoupling AST traversal from model building, the system now employs an extensible visitor pattern.

* The `generate mocks` command does now support additional annotations; use `//parsley:mock` and `//parsley:ignore` to gain full control over how mock generation is handled, while keeping the default behavior of including all interfaces. If `//parsley:mock` is present, it takes precedence, meaning all interfaces are excluded by default, and only those explicitly marked with `//parsley:mock` are included.


## [v0.9.2] - 2024-09-09

### Added

* Added a `version` command to display the current version of the Parsley CLI application. The command can also check for new versions by querying the latest release information from GitHub, notifying users if an update is available and providing instructions on how to update.

### Fixed

* Fixes matching of expected arguments in `mock_verify.go`
* Adds the `FormatType` template function and fixes the type formatting for array types in default mock functions.
  

## [v0.9.1] - 2024-09-09

### Changed

* Added improved testability for the `TemplateModelBuilder` by refactoring its constructor function to accept an `AstFileAccessor` function instead of a filename. This allows for greater flexibility in testing, as the AST can now be sourced either from a file or directly from a string, making it easier to test different code inputs without relying on file I/O.

### Fixed

* Fixed an issue in the `type_model_builder.go` module where parameters and result fields of type array were not correctly handled. This update ensures that array types are properly represented in the generated template models, allowing for accurate code generation in cases involving arrays.


## [v0.9.0] - 2024-09-08

Starting with this release, the project's license has been changed from AGPLv3 to Apache License 2.0. The move to the Apache 2.0 license reflects my desire to make the library more accessible and easier to adopt, especially in commercial and proprietary projects.

### Added

* Adds the `generate mocks` CLI command that can generate configurable mock implemetations from interface types.

### Changed

* Several refactorings to the internal `generator` package with improvements to error handling and extensibility.

* Adds the `generic_generator.go` module, integrating generator templates, output file configuration, and template execution. The initial implementation resided in the `generate_proxy_command.go` module. By pulling variables and control structures from parameters, the generator command logic could be moved to the (internal) `generator` package, allowing the logic to be reused by other code-file generator commands. Adding other template-based generators based on (interface) type models can be achieved with less effort.

* Removes methods from the generator type model - uses a function map instead.

* The generic code generator now formats generated code in canonical go fmt style.


### Fixed

* Fixes generator command short description texts


## [v0.8.3] - 2024-09-01

### Fixed

* Allows registration of (immutable) struct dependencies


## [v0.8.2] - 2024-08-29

### Fixed

* Minor changes to the bootstrap generator templates


## [v0.8.1] - 2024-08-29

Parsley is extended by the `parsley-cli` utility application, which is the foundation for new library features that cannot be implemented on top of reflection. Support for proxy and/or decorator types is better integrated via a code generator approach.

### Added

* Adds the `parsley-cli` application that adds code generation capabilities. 

* The `init` command bootstraps a new Parsley application (a `main.go` and an `application.go` file providing the bare minimum to kick-start a dependency injection-enabled app).

* The `generate proxy` command generates extensible proxy types by `MethodInterceptor` objects, which can function as proxies or decorator objects.
  

## [v0.7.1] - 2024-08-10

This version addresses issues with resolving and injecting services as lists.

### Added

* Adds the `RegisterList[T]` method to enable the resolver to inject lists of services. While resolving lists of a specific service type was already possible by the `ResolveRequiredServices[T]` method, the consumption of arrays in constructor functions requires an explicit registration. The list registration can be mixed with named service registrations.

### Changed

* Changes the key-type used to register and lookup service registrations (uses `ServiceKey` instead of `reflect.Type`). 

* Adds `fmt.Stringer` implementations to registration types to improve the debugging experience. It also fixes the handling of types reflected from anonymous functions.

* Extracts some registry and resolver errors.


## [v0.7.0] - 2024-08-05

### Added

* Adds the `RegisterLazy[T]` method to register lazy service factories. Use the type `Lazy[T]` to consume a lazy service dependency and call the `Value() T` method on the lazy factory to request the actual service instance. The factory will create the service instance upon the first request, cache it, and return it for subsequent calls using the `Value` method.

## [v0.6.1] - 2024-07-30

### Changed

* Registers named services as transient services to resolve them also as a list of services (like services without a name). Changes the `createResolverRegistryAccessor` method so temporary registrations are selected first (and shadow permanent registrations). This behavior can also be leveraged in `ResolverOptionsFunc` to shadow other registrations when resolving instances via `ResolveWithOptions.`


## [v0.6.0] - 2024-07-26

### Added 

* Adds the `Activate[T]` method which can resolve an instance from an unregistered activator func.

* Allows registration and activation of pointer types (to not enforce usage of interfaces as abstractions).

* Adds the `RegisterNamed[T]` method to register services of the same interface and allow to resolve them by name.

### Changed

* Renames the `ServiceType[T]` method to `MakeServiceType[T]`; a service type represents now the reflected type and its name (which makes debugging and understanding service dependencies much easier).

* Replaces all usages of `reflect.Type` by `ServiceType` in all Parsley interfaces.

* Changes the `IsSame` method of the `ServiceRegistration` type; service registrations of type function are always treated as different service types.

### Fixed

* Fixes a bug in the `detectCircularDependency` function which could make the method get stuck in an infinite loop.


## [v0.5.0] - 2024-07-16

### Added

* The service registry now accepts multiple registrations for the same interface (changes internal data structures to keep track of registrations; see `ServiceRegistrationList`).

* Adds the `ResolveRequiredServices[T]` convenience function to resolve all service instances; `ResolveRequiredService[T]` can resolve a single service but will return an error if service registrations are ambiguous.

### Changed

* Extends the resolver to handle multiple service registrations per interface type. The resolver returns resolved objects as a list. 


## [v0.4.0] - 2024-07-13

### Added

* Support for factory functions to create instances of services based on input parameters provided at runtime

### Changed

* Reorganizes the whole package structure; adds sub-packages for `registration` and `resolving`. A bunch of types that support the inner functionality of the package have been moved to `internal.`.

* Integration tests are moved to the `internal` package.


## [v0.3.0] - 2024-07-12

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


## [v0.2.0] - 2024-07-11

### Added

* The resolver can now detect circular dependencies.

* Adds helpers to register services with a certain lifetime scope.

### Changed

* The registry rejects non-interface types.

### Fixed

* Fixes error wrapping in custom error types.

* Improves error handling for service registry and resolver.


## [v0.1.0] - 2024-07-10

### Added

* Adds a service registry; the registry can map interfaces to implementation types via constructor functions.

* Assign lifetime behaviour to services (singleton, scoped, or transient).

* Adds a basic resolver (container) service.
