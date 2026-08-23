# CLAUDE.md — aha-go

Repo-specific guidance for Claude Code. Also read the org-level
`~/go/src/github.com/grokify/.github/CLAUDE.md` for grokify-wide conventions.

## What this repo is

A Go SDK for Aha! (aha.io) with two client layers:

- **REST** (`internal/api`, ogen-generated) — most CRUD, wrapped by
  hand-written types in `feature.go`, `initiative.go`, `release.go`, etc.
- **GraphQL** (`graphql/generated`, genqlient-generated) — used where REST
  can't reach (e.g. custom field writes on Release), wrapped by
  `graphql/custom_fields.go` and similar.

Prefer the REST wrapper for a given field/entity unless GraphQL is the only
path Aha exposes for it (check `openapi/aha.yaml` and Aha's docs at
https://www.aha.io/api first).

`omniroadmap/` is a third, narrower layer: it adapts `*aha.Client` to
`omniroadmap-core`'s `provider.Provider` interface for external roadmap-
aggregation consumers (e.g. `plexusone/dashforge`). It only wraps existing
REST wrapper methods — it never talks to `internal/api` directly, and
changes to it don't touch the OpenAPI spec.

## The OpenAPI spec is hand-maintained — on purpose

**Aha does not publish an OpenAPI spec.** `openapi/aha.yaml` is authored and
verified by hand against Aha's actual documented/observed API behavior, not
vendored or fetched from anywhere. Producing (and proving correct) this spec
is itself part of this SDK's value — treat inaccuracies in it as bugs in our
own research, not as an upstream sync problem.

Workflow for adding/changing a field or endpoint:

1. Verify the field against Aha's official docs (https://www.aha.io/api) —
   don't guess types or nullability.
2. Edit `openapi/aha.yaml`. Follow existing conventions:
   - No enums for account/workflow-specific string fields (e.g.
     `workflow_status`, `progress_source`) — plain `type: string`, validated
     server-side only. Only `User.role` (response-only) is a real enum.
   - `nullable: true` (no `required:`) on a date/number field → ogen
     generates `OptNil*`; omit `nullable` for a plain `Opt*`.
3. Run `make generate` (`./generate.sh`) to regenerate `internal/api`. Never
   hand-edit generated files.
4. Diff the regenerated files and confirm the change is scoped to exactly
   the fields you added — no incidental changes elsewhere in the spec.
5. Wire the new generated `Opt*` fields into the relevant `With*`
   functional-option wrapper in `feature.go`/`initiative.go`/`release.go`
   (etc.), following the existing pattern in that file.
6. Add an `httptest`-based test using the shared `newTestClient` helper
   (`idea_test.go`) — capture the request body into `map[string]any` and
   assert on the JSON keys written, or decode a fixture response and assert
   on the populated struct fields.

## Testing conventions

- Unit tests only, `httptest`-based — no live Aha credentials required.
- Test file per entity (`feature_test.go`, `release_test.go`, ...).
- Use the shared `newUpdateCaptureClient` helper (`idea_test.go`) for
  Update* tests that need to assert on exactly what was sent — it decodes
  the request body into `map[string]any` and responds with a given
  fixture, replacing what used to be a duplicated inline pattern per file.
