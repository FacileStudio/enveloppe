# enveloppe — Configuration

A types-only contract reads no environment variables and has no runtime settings. What it does
have is a set of compile-time and serialization constraints that a consumer has to match, and
those are what this page documents.

## Environment variables

None. Neither `go/enveloppe.go` nor `ts/src/index.ts` reads the environment, opens a file, or
makes a network call. If a consumer needs configuration to *send* an event, that belongs to
[pool](https://github.com/FacileStudio/pool) — see its `nook.yaml` reference.

## Go consumers

| Setting | Value | Why |
|---|---|---|
| Module path | `github.com/FacileStudio/enveloppe/go` | The `/go` suffix is part of the path; the module lives in a subdirectory of the repo |
| Minimum Go | 1.24.0 | Declared in `go/go.mod` |
| Dependencies | none | The module requires nothing |

There are **no semver tags** on this repo, so `go get` resolves to a pseudo-version of a
commit on `main`:

```
require github.com/FacileStudio/enveloppe/go v0.0.0-20260804090730-02b0f4b20c6f
```

Sablier and Mycelium are both pinned to that commit. Updating is deliberate:
`go get github.com/FacileStudio/enveloppe/go@<commit>`, then read the diff — a contract change
is never a routine bump.

Go 1.24 is required because `Event[T]` uses generics, and the module declares that floor
regardless. It imposes nothing else on the consumer.

## TypeScript consumers

`ts/tsconfig.json` compiles the package with:

| Option | Value | Consequence for consumers |
|---|---|---|
| `target` | `ES2020` | |
| `module` / `moduleResolution` | `NodeNext` | Consumers must resolve ESM. `package.json` sets `"type": "module"` and exports `./dist/index.js` |
| `declaration` | `true` | `dist/index.d.ts` is the actual product |
| `strict` | `true` | Optional fields arrive as `string \| undefined` under a consumer's own `strict` |
| `skipLibCheck` | `true` | |

`package.json` exposes a single `.` export mapping `types` to `./dist/index.d.ts` and `import`
to `./dist/index.js`, with `"files": ["dist"]`. There is no CommonJS entry point — a CJS
consumer cannot import this package as published.

The distribution caveat is in [development.md](development.md#distribution): the `#ts` branch
the install command expects does not exist, so no TypeScript consumer can currently install it.

## Serialization rules

These are the constraints that make the two languages produce byte-compatible JSON. Match them
when adding a field.

| Intent | Go | TypeScript | On the wire |
|---|---|---|---|
| Required | `string` with a plain tag | `field: string` | Always present |
| Nullable | `*string`, **no** `omitempty` | `field: string \| null` | Always present, `null` when empty |
| Optional | `string` with `,omitempty` | `field?: string` | Absent when unset |

Getting nullable and optional the wrong way round changes the wire format silently — a
receiver that checks `"stopped_at" in payload` behaves differently from one that checks
`payload.stopped_at !== null`. Copy the shape of an existing field rather than deciding fresh.

Timestamps are plain strings in both languages. The contract does not parse or validate them;
RFC 3339 in UTC is the convention across the suite, and the receiver is responsible for
parsing.

Numbers: `amount` is `float64` in Go, and `tokens_in` / `tokens_out` are `int64`. JSON has one
number type, so a TypeScript consumer sees `number` for all three and inherits the usual
double-precision limit — irrelevant for token counts, worth remembering for money. The contract
does not carry currency minor units; `amount` and `currency` are separate free-form fields.

## Envelope version

`version` is `1` and has never moved. Read it from the exported constant —
`enveloppe.EventVersion` or `FACILE_EVENT_VERSION` — rather than writing a literal, so a bump
propagates through recompilation. Nothing in this repo validates it; consumers that want to
reject a future version have to check it themselves.
