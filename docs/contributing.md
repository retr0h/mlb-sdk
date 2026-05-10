# Contributing to mlb-sdk

Thanks for your interest in contributing.

## Before you start

- **Read [AI_POLICY.md](../AI_POLICY.md)** — disclose AI assistance, ensure you
  understand any code you submit.
- **Check existing work** — search open issues and PRs to avoid duplicating
  effort.
- **Start small** — focused PRs are easier to review than sweeping changes.

## Making changes

### Code style

- Run `gofmt -l .` and `go vet ./...` before pushing.
- Tests: `go test ./...`. Add tests for any non-trivial behavior.
- Idiomatic Go on the public surface — see
  [development.md](development.md#public-surface-authoring).

### OpenAPI spec changes

- Authoring rules in [development.md](development.md#openapi-spec-authoring) —
  named components only, no inline nested objects.
- After editing `api/openapi.yaml`, regenerate with `go generate ./internal/gen`
  and commit the regenerated `internal/gen/client.gen.go`.

### Documentation

- Update docs alongside code changes when public behavior changes.

## Submitting a PR

1. Create a feature branch from `main` (`type/short-description` — `feat/`,
   `fix/`, `docs/`, `refactor/`, `chore/`).
2. Commit messages:
   [Conventional Commits](https://www.conventionalcommits.org/), see
   [development.md](development.md#commit-messages).
3. PR description: what changed, why, and any follow-ups.
4. Open as draft if you want early feedback before final review.
5. One logical change per PR — split unrelated changes.

## Adding endpoints

1. Add the path + named schemas to `api/openapi.yaml`.
2. `go generate ./internal/gen` — verify the new operation appears in
   `internal/gen/client.gen.go`.
3. Add public types in `pkg/mlb/<resource>_types.go`.
4. Add the method on `*Client` in `pkg/mlb/<resource>.go`, plus any helper
   methods on the response types that hide upstream awkwardness.
5. Add tests in `pkg/mlb/<resource>_test.go` — prefer the `httptest.Server` +
   `WithBaseURL` pattern from `boxscore_test.go`.
