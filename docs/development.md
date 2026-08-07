# enveloppe — Development

Building both halves, how each one reaches its consumers, and the rules for changing a contract
that other repos compile against.

## Prerequisites

- **Bun**, for the TypeScript half. `ts/bun.lock` is committed.
- **Go 1.24**, for the Go half. Nothing to install beyond the toolchain — the module has no
  dependencies.

## Build

```sh
cd ts && bun install && bun run build
```

```sh
cd go && go build ./...
```

`bun run build` is `tsc`. It emits `ts/dist/index.d.ts` — the real product — and
`ts/dist/index.js`, which is one line, because `FACILE_EVENT_VERSION` is the only value in the
package that survives type erasure.

There is **no test suite, no linter, no CI and no dev server**, and for a types-only contract
that is defensible: the Go compiler and `tsc` are the entire check. What they cannot catch is
the two halves drifting apart, which is the actual risk here — see below.

## Committed build output

`ts/dist/` is committed on purpose. A `github:` dependency installs the repository as-is and
runs no build step, so an uncommitted `dist/` would ship a package with no types at all.

**Run `bun run build` and commit `ts/dist/` in the same commit as any `ts/src/` change.** A
source-only commit ships nothing and does so silently, because the stale `dist/` still
resolves.

## Distribution

Neither half goes through a registry.

**Go — a commit on `main`.** The module path `github.com/FacileStudio/enveloppe/go` resolves the
`go/` subdirectory from the default branch. There are no semver tags, so consumers pin a
pseudo-version; Sablier and Jardin are both on `v0.0.0-20260804090730-02b0f4b20c6f`. A push to
`main` does not move them — Go consumers update deliberately.

**TypeScript — currently not installable.** `package.json` and the historical install
instructions expect `bun add github:FacileStudio/enveloppe#ts`, mirroring the arrangement in
[antenne-client](https://github.com/FacileStudio/antenne-client), where a `ts` branch holds the contents of `ts/`
at the repository root so Bun can find a `package.json` there.

**That branch has never been published.** `FacileStudio/enveloppe` has exactly one branch,
`main`, and `main` has no `package.json` at its root — so the `#ts` reference cannot resolve and
a rootless install would find no package either. There are correspondingly **zero TypeScript
consumers** of `@facile/enveloppe` in the suite; every importer today is Go.

Nothing here is broken in production, because nothing depends on it. But the TypeScript half is
maintained code that no one can install, so treat any claim that a TS app "uses enveloppe" as
false until the branch exists. Publishing it means mirroring `ts/` to a `ts` branch the way
`pool` does, and once it exists a push to that branch is an immediate release to every
consumer — `#ts` is a moving reference, not a version.

## Changing the contract

Other repos compile against these types. The Go side is pinned per-commit, so a change lands
only when a consumer opts in; there is no such brake on the TypeScript side once `ts` exists.

**Non-breaking**, add freely: a new optional field (`omitempty` / `?`), a new payload type, a
new `app` value.

**Breaking**, and it needs the `version` constant bumped plus a coordinated update of every
consumer: renaming a JSON key, removing a field, making an optional field required, changing a
field's type or nullability, removing an enum value.

**Ambiguous, and worse than either:** adding a new `action` or `object` value. It is additive
in the emitter and a silent no-op in every receiver that has not been updated — the event
arrives, matches no branch, and is dropped without an error. Ship the receivers first.

Whatever you change, change **both languages in the same commit**. Nothing verifies that they
agree; the two files are kept in sync by hand and by review. Compare the JSON tags in
`go/enveloppe.go` against the field names in `ts/src/index.ts` line by line — they are the
contract, and a Go field renamed without its tag updated changes nothing on the wire while a
tag changed without the TypeScript field breaks everything.

The commit history is the changelog — there is no `CHANGELOG.md`. Messages carry the reason:
`feat: agent_session object + Jardin app for agent time tracking`,
`fix: restore actor_email field on Task and update docs to match code`. Keep that habit; it is
the only record of why a field exists.

## Known outstanding fix

`AgentSession.UserEmail` is a required field while every other email in the contract is
optional. Making it `omitempty` in Go and `?` in TypeScript is the consistent shape. It is
non-breaking for Jardin, which always supplies one.
