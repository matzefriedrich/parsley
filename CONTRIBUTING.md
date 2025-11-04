# Contributing to Parsley

I appreciate your interest in Parsley. Contributions are appreciated, and this document helps ensure that every change supports the project’s long-term quality and goals.


## Project philosophy

Parsley is designed to **help developers manage complexity**, not pretend it doesn’t exist. Go encourages straightforward code, but modern applications often require stronger structure. Contributions should improve clarity, navigability, and maintainability of large systems - even when that means introducing new abstractions.

Simplicity is valuable, but only when it reduces cognitive load, not when it hides real-world complexity.


## Principles

* Focus on improving project clarity and evolvability
* Favor structured and well-designed abstractions over ad-hoc minimalism
* Contributions should feel “at home” in the codebase
* If something adds developer friction, reconsider the approach


## Contribution process

### Discuss first when changing behavior

Please open an issue for any non-trivial feature or changes to public behavior, describing:

* The problem being solved
* Design options considered
* Expected impact on users and future contributors

This avoids wasted effort if an idea doesn’t align with the project direction.


### Development workflow

The project maintains a clean and readable Git history to support efficient review and long-term maintainability.
Create branches following one of these patterns:

```
feature/<short-summary>
fix/<issue-number>
```

* Please ensure that commits are logical and easy to read. Aim for a linear commit history.
* During development, commits may be exploratory or incremental, but before review:
  * Use `git rebase` and `--autosquash` with `fixup!` / `squash!` commits
  * Consolidate small or sequential commits into meaningful units
  * Each commit should represent a clear step in the change, not partial work


#### Tooling requirements

The goal is to keep the development setup lightweight and accessible. Use the **standard Go SDK** tooling. Please don't just introduce new linting frameworks or build dependencies unless they have been previously discussed and justified.


## Changelog requirements

Every pull request must update the `CHANGELOG.md` file:

* Add an entry under the **Unreleased** section
* Provide a concise, meaningful description, and include a link to the pull request

Small changes count - everything visible to users should be noted.


## Code standards

The project aims to provide a framework that supports developers in creating flexible and extensible project structures for Go applications, ensuring they can adapt to evolving requirements and changing environmental constraints.

* Use existing patterns and terminology as a guide
* Abstractions are welcome when they **clarify structure**
* Keep dependencies minimal and intentional

### Error-handling and observability

Parsley does not include internal logging, since logging strategies vary across applications. Instead, contributions must focus on structured error handling that allows applications using Parsley to integrate their own logging and observability tools cleanly.

* Errors should provide actionable context
* Avoid returning plain strings when a meaningful error type is appropriate
* Do not introduce debug prints or hidden side-effects
* Let the consuming application decide how and when to log or report issues
* When in doubt, prefer returning richer error information over emitting logs

The goal is to enable strong diagnostics where they belong: in the application that owns the operational concerns.


### Testing requirements

* Tests must pass reliably on any supported platform
* Include edge cases and error conditions
* Validate concurrency behavior when relevant

If CI breaks, please update the PR to restore a passing state before review continues.


## Pull-request expectations

Always rebase onto the latest main branch if the PR falls behind. Keeping a PR up to date avoids merge conflicts at the end and ensures CI verifies the correct state.

A good PR includes:

* Clear explanation of the change and motivation
* Links to any related issues
* Relevant tests for new or changed behavior
* An Unreleased CHANGELOG entry

If revisions are requested, they are aimed at making the contribution stronger - not at discouraging you.


## Review and collaboration

Reviews focus on:

* Technical correctness
* Long-term maintainability
* Architectural alignment with project goals

If a PR stalls for an extended period without response, it may be closed, but contributions are always welcome again after revision.
