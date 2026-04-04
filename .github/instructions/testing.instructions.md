---
description: "Use when writing, reviewing, or fixing Go tests in this repository. Covers AAA structure, fixtures, naming, and module-level validation."
applyTo: "**/*_test.go"
---

# Testing Instructions — go-redux (Go)

## 1) Structure and Naming

- Every test must use explicit `// Arrange`, `// Act`, `// Assert` comments.
- A single test must have exactly one `Act`.
- If a scenario needs multiple operations, split it into multiple tests and share setup through fixtures or helpers.
- Prefer names like `TestStoreDispatch_ShouldPublishStateChange_WhenActionChangesState`.

## 2) AAA and Fixtures

- Keep `Arrange` small. If setup grows beyond a few lines, extract helpers such as `newTestStore`, `newMutedLogger`, or `newCounterHandler`.
- Fixtures must build real domain objects and real handlers when practical.
- Use table-driven tests only when each case still represents one clear `Act`.

## 3) Assertions

- Use the Go standard library (`testing`, `errors`, `reflect`, `slices`, `cmp` only if already present).
- Assert the specific behavior that matters to the scenario.
- Do not hide multiple assertions behind loops unless the loop is the scenario itself.

## 4) Error and Panic Testing

- Prefer validating returned errors over recovering from panics.
- If panic behavior must be tested, the panic check is the single `Act`.
- Runtime paths should usually return errors; configuration/setup paths may panic for programmer misuse.

## 5) Isolation and Determinism

- No shared mutable globals between tests.
- Avoid timers and sleeps when a channel, wait group, or synchronous API can prove the behavior directly.
- Mute loggers in tests unless the output is part of the assertion.

## 6) Module Coverage

- Root module changes should add or update tests in the root module.
- Example modules are contracts too: if an API change affects examples, run the affected example module tests.
- CI must cover the root module and every nested example module with its own `go.mod`.
