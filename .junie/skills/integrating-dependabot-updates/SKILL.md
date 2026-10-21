---
name: integrating-dependabot-updates
description: Use when integrating a Dependabot dependency-update pull request in this Go repository, especially when rebasing, resolving conflicts, tidying modules, applying lint fixes, or updating the changelog.
---

# Integrating Dependabot Updates

Use this workflow when asked to integrate a Dependabot pull request. Keep the PR's dependency changes isolated, validate automated edits, and record user-visible changes without rewriting released notes.

## Workflow

1. Inspect the PR state, base and head branches, commits, files, and checks. Confirm the repository and PR identity.
2. Check the current worktree and branch status before fetching. Preserve existing edits and untracked files. When the default worktree is on another branch or contains unrelated work, use a separate PR worktree.
3. Fetch the PR head and target branch. Verify the fetched head belongs to the intended repository and matches the PR branch before checking it out; then rebase it onto the updated target branch. Do not force-push unless explicitly requested.
4. Resolve conflicts one file at a time using the target and PR intent. Do not accept one side wholesale by default. Review the resulting diff and confirm all intended PR changes remain.
5. In the PR worktree, run `go mod tidy`, then review `go.mod` and `go.sum` for unexpected dependency changes. Run `make lint-fix` and inspect every resulting edit.
6. Update `CHANGELOG.md` under `## [unreleased]` with concise notes that describe the actual developer-facing change. Keep previously released sections unchanged and avoid duplicating existing entries.
7. Run relevant tests and checks after all edits. Inspect the complete diff, `git status`, and commit ancestry; report any unresolved failures or follow-up actions. Do not commit or push unless requested.

## Conflict and Change Review

- A successful rebase does not prove the intended behavior survived; inspect the rebase result and run tests.
- Treat `go mod tidy` and lint autofixes as edits to review, not as automatically correct or disposable output.
- Keep dependency-only bumps concise in the changelog; explain behavior or compatibility impact when applicable.
- Never discard pre-existing user changes or remove unrelated untracked content.

## Common Mistakes

- Running commands from a different worktree than the PR branch.
- Resolving conflicts by choosing “ours” or “theirs” for an entire file without understanding both sides.
- Leaving PR notes in a released changelog section or copying a release entry twice.
- Claiming tests, lint, or a rebase passed without checking the command output and final repository state.