# enveloppe — Architecture

What the envelope is, where it sits between apps, who owns each object type, and the identity
model the whole contract rests on.

## Runtime topology

`enveloppe` has no runtime. It is compiled into both sides of every cross-app conversation and
disappears at build time — the only thing that exists at runtime is the JSON it describes.

```
Mycelium (Go) ──┐ enveloppe.Event[AgentSession]
Sablier (Go)  ──┤ marshaled to JSON
Opus (TS)     ──┘        │
                         ▼
                      pool client  ──▶ WS ──▶ Nook Pool ──▶ WS ──▶ pool client
                                                                        │
                         ┌──────────────────────────────────────────────┘
                         ▼
               FacileEvent<T> / Event[T]  unmarshaled by the receiving app
```

Three repos, one contract:

| Repo | Role |
|---|---|
| `enveloppe` | The type contract. What an event *is* |
| `pool` | The client. How an event *travels* |
| `Nook` | The broker. Where every event *goes* |

No app imports another app's types. That is the point: adding an app means agreeing with
`enveloppe`, not with six other codebases.

## The envelope

Every message on the wire is an envelope wrapping exactly one domain object:

```json
{
  "version": 1,
  "app": "Mycelium",
  "object": "agent_session",
  "action": "created",
  "facile_id": "fac_abc123",
  "payload": { "facile_id": "fac_abc123", "...": "..." },
  "timestamp": "2026-08-05T14:03:11Z",
  "idempotency_key": "mycelium_agent_session_created_fac_abc123"
}
```

`app`, `object` and `action` are closed sets — see [api.md](api.md). `facile_id` is duplicated
at the envelope level so a consumer can route or deduplicate without unmarshaling the payload
into a concrete type.

`idempotency_key` is what makes the whole thing survivable. `pool` keeps its replay offsets in
memory only, so a consumer that restarts replays from offset zero on every channel and will see
events it has already applied. The key is the receiver's defense; nothing in this repo
generates or validates it. The convention visible in the suite is
`<app>_<object>_<action>_<facile_id>`.

`timestamp` is a string in both languages — RFC 3339 by convention, unparsed by the contract.
Note that `pool` puts its own numeric `timestamp` (epoch milliseconds) on the transport message
around this one; they are different fields at different layers and both reach a handler.

## Channels

The contract does not name channels. Consumers derive them as `<object>.<action>` —
`project.created`, `task.updated`, `agent_session.created` — as declared in each app's
`nook.yaml`. Nothing here enforces that mapping, so an app can emit on any channel string
`pool` will accept.

## Object ownership

Ownership is convention, not enforcement — the contract cannot stop an app from emitting an
object it does not own. Only two repos import this module today (Sablier and Mycelium, both Go),
so most of the enum is declared ahead of any producer.

| Object | Emitted by | Notes |
|---|---|---|
| `project` | Sablier, Opus | The one genuinely bidirectional pair in the suite today |
| `task` | Sablier, Opus | Same pair |
| `time_entry` | Sablier | Mycelium **never** emits this, though it produces the raw activity |
| `agent_session` | Mycelium | One event per sealed block of AI-agent activity |
| `monitor` | Sonde | One event per uptime probe state change; Antenne consumes it for alerting |
| `user`, `invoice`, `document` | nobody yet | Declared in the contract; no importer emits them |

Likewise, `Ardoise`, `Plume`, `Glouton` and `Vision` are valid values of `app` but do not
import this module. The enum is the agreed vocabulary, not a record of what is wired.

The `agent_session` → `time_entry` split is worth understanding, because it is the model for
every future pair. Mycelium emits `agent_session` — machine, agent, branch, `user_email`,
start/stop, token totals — for a sealed block of agent activity. Sablier listens and
materializes those as time entries. Mycelium does not emit `time_entry` because that object
belongs to Sablier, and an object with two producers has no owner at all.

## Identity

**`facile_id` is the canonical cross-app join key.** Every domain struct leads with it. Email is
a display attribute, not a key.

That decision is load-bearing, and the reason for it is that **three apps physically cannot
supply an email**:

- **Capsule** — has no user table by design; it is an end-to-end encrypted paste service
- **Perception** — stores `actor_email_hash`, deliberately one-way
- **Nook** — keys on apps, not on users

Keying the contract on email would have excluded all three permanently. `facile_id` also
survives the move to a self-hosted IdP, where Authentik issues the stable subject and email
becomes just another claim that can change.

The email fields that remain are optional accordingly: `Task.actor_email` and
`TimeEntry.user_email` are `omitempty` in Go and `?`-optional in TypeScript, so an app with no
email to offer emits a valid event by leaving them out.

**Known gap:** `AgentSession.user_email` is **not** optional — it is a required `string` in
both languages. Only Mycelium emits `agent_session` and Mycelium does have emails, so nothing is
broken today, but it is inconsistent with the rule above and blocks any future producer that
cannot supply one.

## Nullability

`description`, `icon` and `stopped_at` are the only genuinely nullable fields: `string | null`
in TypeScript, `*string` in Go, with no `omitempty`. They serialize as an explicit `null`,
which is the distinction the contract wants — "this timer has not stopped" is not the same as
"this event did not mention stopping".

Optional fields go the other way: `?` in TypeScript, `omitempty` in Go, absent from the JSON
entirely when unset. Getting these two mixed up changes the wire format, so match the existing
fields rather than guessing.

## Versioning the contract

`version` on the envelope, `EventVersion` / `FACILE_EVENT_VERSION` as the constant, currently
`1`. It has never been bumped. Nothing in this repo reads it — it exists so a consumer can
branch on it once a breaking change happens. See
[development.md](development.md#changing-the-contract) for what counts as breaking.
