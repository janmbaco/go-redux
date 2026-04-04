# go-redux Repository Instructions

These instructions apply to the whole repository. Keep code, tests, examples, and documentation aligned with the exported v2 API.

## Architecture

- The public library surface lives in the root module `github.com/janmbaco/go-redux/v2`.
- Keep the root package focused on ergonomic exports and store composition.
- `actions/`, `handlers/`, `events/`, `reducers/`, and `ioc/` contain the implementation details behind the public API.
- Examples are independent Go modules. If the root API changes, update every affected example module in the same change.
- Do not introduce framework-heavy abstractions when a small Go type or function is enough.

## Code Style

- Write exported Go doc comments in English.
- Prefer explicit, typed APIs over `map[string]any` unless the API is intentionally dynamic.
- Return errors instead of panicking in runtime paths. Panics are acceptable only for programmer errors during configuration.
- Keep concurrency explicit. Protect shared mutable state with mutexes or other synchronization primitives.
- Preserve backward compatibility inside the v2 line unless a breaking change is intentional and documented.

## Public API Discipline

- README examples must compile against the code in the repo.
- `README.md`, `CHANGELOG.md`, and `docs/**` must describe the current API, not a planned API.
- The root library module path is `github.com/janmbaco/go-redux/v2`; do not add new root documentation or root-library examples using the v1 path.
- If a helper is presented as part of the root package in docs, either export it from the root package or update the docs immediately.

## Build and Validation

- Root module: `go test ./...`
- Example modules: run `go test ./...` inside each module under `examples/**/go.mod`
- Keep code formatted with `gofmt`.
- Prefer standard-library tests and helpers unless an additional dependency is clearly justified.

## Linked References

- Test rules: [.github/instructions/testing.instructions.md](instructions/testing.instructions.md)
- Documentation rules: [.github/instructions/docs.instructions.md](instructions/docs.instructions.md)
- Repository overview: [README.md](../README.md)
