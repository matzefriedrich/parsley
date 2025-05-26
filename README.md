[![CI](https://github.com/matzefriedrich/parsley/actions/workflows/go.yml/badge.svg)](https://github.com/matzefriedrich/parsley/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/matzefriedrich/parsley/badge.svg)](https://coveralls.io/github/matzefriedrich/parsley)
[![Go Reference](https://pkg.go.dev/badge/github.com/matzefriedrich/parsley.svg)](https://pkg.go.dev/github.com/matzefriedrich/parsley)
[![Go Report Card](https://goreportcard.com/badge/github.com/matzefriedrich/parsley)](https://goreportcard.com/report/github.com/matzefriedrich/parsley)
![License](https://img.shields.io/github/license/matzefriedrich/parsley)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matzefriedrich/parsley)
![GitHub Release](https://img.shields.io/github/v/release/matzefriedrich/parsley?include_prereleases)

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B48327%2Fgithub.com%2Fmatzefriedrich%2Fparsley.svg?type=shield&issueType=license)](https://app.fossa.com/projects/custom%2B48327%2Fgithub.com%2Fmatzefriedrich%2Fparsley?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B48327%2Fgithub.com%2Fmatzefriedrich%2Fparsley.svg?type=shield&issueType=security)](https://app.fossa.com/projects/custom%2B48327%2Fgithub.com%2Fmatzefriedrich%2Fparsley?ref=badge_shield&issueType=security)

## What is Parsley?

Parsley is a reflection-based dependency injection (DI) package for Go that streamlines dependency management through automated lifetime management, type-safe registration, service resolution and activation, proxy generation, and method interception. It leverages reflection only at startup, ensuring runtime efficiency, and supports features like mocking and testing integration for better development workflows.

### Why use dependency injection in Go?

Dependency injection in Go typically relies on constructor functions, which are idiomatic and foundational for creating and managing dependencies. While this approach is powerful, it can lead to boilerplate code as projects become more complex, requiring explicit instantiation and manual wiring. Addionally, more sophisticated ramp-up code is needed, if not all instances shall be created at application start. Parsley eliminates much of this repetitive code by automating the creation and wiring of dependencies, enabling developers to focus on application logic. It aims to enhance modularity and make managing dependencies easier.

## Key Features

- **Type Registration**: Supports constructor functions, lifetime management (singleton, scoped, transient), and safe casts.  
- **Modular and Lazy Loading**: Register types as modules, resolve dependencies on-demand, and inject lazily with `Lazy[T]`.  
- **Registration Validation**: Ensures early detection of missing registrations or circular dependencies.
- **Advanced Registrations**: Register multiple implementations for the same interface using named services or lists.  
- **Proxy and Mock Support**: Generate extensible proxy types and configurable mocks to streamline testing workflows.  

For a complete overview of features and capabilities, refer to the [Parsley documentation](https://matzefriedrich.github.io/parsley-docs/).  


## Usage

### Install Parsley

Add Parsley to your project:

```sh
go get github.com/matzefriedrich/parsley
```

For advanced features like proxy generation or mock creation, install the `parsley-cli` utility:

```sh
go install github.com/matzefriedrich/parsley/cmd/parsley-cli
```

--- 

Copyright 2024 - 2025 by Matthias Friedrich
