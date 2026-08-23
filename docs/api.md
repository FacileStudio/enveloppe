# enveloppe — API

Every type in the contract, every field, its JSON key, whether it is required, and how the two
implementations differ. Taken from `go/enveloppe.go` and `ts/src/index.ts`; where they
disagree, that is called out explicitly.

"Required" below means the field is always present in the serialized JSON. A field marked
optional is `omitempty` in Go and `?` in TypeScript, and is **absent** when unset. A field
marked nullable is always present and carries `null` when empty.

## Enumerations

### `app` — `enveloppe.App` / `FacileApp`

Eight values. Go exports one constant each; TypeScript is a string-literal union.

| Value | Go constant |
|---|---|
| `Sablier` | `AppSablier` |
| `Opus` | `AppOpus` |
| `Ardoise` | `AppArdoise` |
| `Plume` | `AppPlume` |
| `Glouton` | `AppGlouton` |
| `Vision` | `AppVision` |
| `Mycelium` | `AppMycelium` |
| `Sonde` | `AppSonde` |

Capitalized exactly as shown — `Sablier`, not `sablier`. Both halves agree.

### `action` — `enveloppe.Action` / `FacileAction`

| Value | Go constant |
|---|---|
| `created` | `ActionCreated` |
| `updated` | `ActionUpdated` |
| `deleted` | `ActionDeleted` |

### `object` — `enveloppe.ObjectType` / `FacileObjectType`

Eight values, one per payload type below.

| Value | Go constant | Payload type |
|---|---|---|
| `project` | `ObjectProject` | `Project` / `FacileProject` |
| `task` | `ObjectTask` | `Task` / `FacileTask` |
| `user` | `ObjectUser` | `User` / `FacileUser` |
| `time_entry` | `ObjectTimeEntry` | `TimeEntry` / `FacileTimeEntry` |
| `invoice` | `ObjectInvoice` | `Invoice` / `FacileInvoice` |
| `document` | `ObjectDocument` | `Document` / `FacileDocument` |
| `agent_session` | `ObjectAgentSession` | `AgentSession` / `FacileAgentSession` |
| `monitor` | `ObjectMonitor` | `Monitor` / `FacileMonitor` |

Nothing in the contract ties an `object` value to its payload type. `Event[T]` and
`FacileEvent<T>` are generic over `T` and will happily carry a mismatched pair.

## The envelope — `Event[T]` / `FacileEvent<T>`

| JSON key | Type | Required | Go field | TS field | What it is |
|---|---|---|---|---|---|
| `version` | number | yes | `Version int` | `version: number` | Contract version; use the exported constant |
| `app` | `App` | yes | `App App` | `app: FacileApp` | Emitting app |
| `object` | `ObjectType` | yes | `Object ObjectType` | `object: FacileObjectType` | Object type carried in `payload` |
| `action` | `Action` | yes | `Action Action` | `action: FacileAction` | What happened to it |
| `facile_id` | string | yes | `FacileID string` | `facile_id: string` | Cross-app id of the object, duplicated from the payload |
| `payload` | `T` | yes | `Payload T` | `payload: T` | The domain object. `T` defaults to `unknown` in TypeScript; no default in Go |
| `timestamp` | string | yes | `Timestamp string` | `timestamp: string` | RFC 3339 by convention; unparsed by the contract |
| `idempotency_key` | string | yes | `IdempotencyKey string` | `idempotency_key: string` | Receiver's deduplication key |

No field is optional or nullable. An envelope with a zero-value Go field still serializes it.

**Version constant:** `enveloppe.EventVersion` (untyped constant) and `FACILE_EVENT_VERSION`
(`const` value `1`). This is the only value the TypeScript build emits at runtime — everything
else is erased, so `ts/dist/index.js` is a single line.

## Payload types

### `project` — `Project` / `FacileProject`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `name` | string | yes | |
| `description` | string \| null | yes, nullable | `*string` in Go |
| `icon` | string \| null | yes, nullable | `*string` in Go |

### `task` — `Task` / `FacileTask`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `project_facile_id` | string | yes | The owning project's `facile_id` |
| `name` | string | yes | |
| `status` | string | yes | Free-form; the contract does not enumerate statuses |
| `actor_email` | string | **no** | `omitempty` / `?`. Absent when the emitter has no email |

### `user` — `User` / `FacileUser`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `email` | string | yes | The one place email is required, because it *is* the object |
| `name` | string | yes | |

### `time_entry` — `TimeEntry` / `FacileTimeEntry`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `task_facile_id` | string | yes | |
| `user_facile_id` | string | yes | |
| `user_email` | string | **no** | `omitempty` / `?` |
| `started_at` | string | yes | |
| `stopped_at` | string \| null | yes, nullable | `null` while the timer is still running |

### `agent_session` — `AgentSession` / `FacileAgentSession`

One event per sealed block of AI-agent activity. Emitted by Mycelium, consumed by Sablier.

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `project` | string | yes | A project **name**, not a `facile_id` — the only such field in the contract |
| `machine` | string | yes | |
| `agent` | string | yes | |
| `branch` | string | **no** | `omitempty` / `?` |
| `user_email` | string | yes | **Known inconsistency** — see below |
| `started_at` | string | yes | |
| `stopped_at` | string | yes | Not nullable: a session is only emitted once sealed |
| `tokens_in` | number | yes | `int64` in Go |
| `tokens_out` | number | yes | `int64` in Go |

`user_email` is required here while `Task.actor_email` and `TimeEntry.user_email` are optional.
Nothing breaks today because Mycelium is the only producer and does have emails, but it
contradicts the identity model in [architecture.md](architecture.md#identity) and blocks any
producer that cannot supply one. Making it `omitempty` is a known outstanding fix.

### `invoice` — `Invoice` / `FacileInvoice`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `client_name` | string | yes | |
| `amount` | number | yes | `float64` in Go |
| `currency` | string | yes | Free-form; no ISO 4217 constraint in the contract |
| `status` | string | yes | Free-form |

### `document` — `Document` / `FacileDocument`

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `facile_id` | string | yes | |
| `title` | string | yes | |
| `status` | string | yes | Free-form |
| `signer_email` | string | yes | |

### `monitor` — `Monitor` / `FacileMonitor`

One event per uptime probe state change. Emitted by Sonde, consumed by Antenne.

| JSON key | Type | Required | Notes |
|---|---|---|---|
| `monitor_id` | string | yes | Sonde's own row ID, alongside the usual cross-app join key |
| `facile_id` | string | yes | |
| `slug` | string | yes | URL-safe identifier, also used by public status pages |
| `name` | string | yes | Human label |
| `status` | string | yes | `"up"` or `"down"`; free-form `string` in Go, literal union in TS |
| `latency_ms` | number | yes | Last probe latency; `0` on a failed probe; `int64` in Go |
| `error` | string \| null | yes, nullable | Last failure message; `null` while the monitor is up |

## Cross-language differences

The two halves are field-for-field identical in JSON. What differs is what the language gives
you around them.

| | Go | TypeScript |
|---|---|---|
| Enum values | Typed constants (`AppMycelium`) | String-literal unions, no constants |
| Generic default | none — `Event[T]` requires `T` | `FacileEvent<T = unknown>` |
| Nullable field | `*string`, dereference before use | `string \| null`, narrow before use |
| Optional field | `omitempty`, zero value is indistinguishable from absent | `?`, `undefined` is distinguishable |
| Runtime output | the constants | one line: `export const FACILE_EVENT_VERSION = 1;` |

The `omitempty` asymmetry is the one that bites. In Go an empty `Task.ActorEmail` is
serialized as absent and unmarshals back to `""`, so a receiver cannot tell "not supplied"
from "supplied as empty". TypeScript preserves that distinction as `undefined`. Do not build
logic on the difference.
