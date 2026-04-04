---
description: "Use when editing README, changelog, docs, or examples. Prevents documentation drift across the v2 API surface."
applyTo: "README.md,CHANGELOG.md,docs/**/*.md,examples/**/README.md,examples/**/go.mod,examples/**/*.go"
---

# Documentation and Example Instructions

- Document the API that exists today, not the API that is planned.
- Keep root library imports on the v2 module path: `github.com/janmbaco/go-redux/v2`.
- Example submodules may use their own module paths, but they must remain internally consistent and still consume the root library through `/v2`.
- If an example is linked from `README.md`, the path must exist in the repository.
- If examples are independent modules, their `go.mod` files and internal imports must stay consistent with the repository layout.
- Migration notes must remain executable enough to guide a real upgrade.
- When renaming exported helpers or changing signatures, update `README.md`, `CHANGELOG.md`, and every example in the same change.
