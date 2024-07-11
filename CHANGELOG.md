# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] - 2024-07-11

### Added

* Service registrations can be bundled in a `ModuleFunc` to register related types as a unit
* Service registry accepts object instances as singleton service registrations
* Adds the `ResolveRequiredService[T]` convenience function that resolves and safe-casts objects
* Registers resolver instance with the registry so that the `Resolver` object can be injected into factory and constructor methods

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
