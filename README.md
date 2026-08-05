# enveloppe

The canonical event envelope for the Facile Suite — the shared type contract every app uses to
talk through Nook, published as `@facile/enveloppe` in TypeScript and `enveloppe` in Go.

Types only. No runtime logic, no validation, no dependencies in either language. This repo is
the contract, not the implementation; [pool](https://github.com/FacileStudio/pool) is the wire
that carries it.

## What it does

- Defines `FacileEvent<T>` / `Event[T]`, the envelope wrapping every cross-app event
- Names the seven apps, three actions and seven object types that may appear on the wire
- Defines the seven domain payload shapes, one per object type
- Pins the envelope version as a constant both languages export
- Keeps `facile_id` as the canonical cross-app join key on every domain object
- Guarantees the two implementations serialize to identical JSON keys

## Stack

| Layer | Tech |
|---|---|
| Runtime | TypeScript 5.8, ES2020/NodeNext, declaration-only output, no dependencies |
| Runtime | Go 1.24, structs and constants with JSON tags, no dependencies |

## Install

Go consumers pin the pseudo-version of a commit — there are no semver tags.

```sh
go get github.com/FacileStudio/enveloppe/go
```

```go
import enveloppe "github.com/FacileStudio/enveloppe/go"

event := enveloppe.Event[enveloppe.AgentSession]{
	Version:        enveloppe.EventVersion,
	App:            enveloppe.AppJardin,
	Object:         enveloppe.ObjectAgentSession,
	Action:         enveloppe.ActionCreated,
	FacileID:       "fac_abc123",
	Payload:        session,
	Timestamp:      time.Now().UTC().Format(time.RFC3339),
	IdempotencyKey: "jardin_agent_session_created_fac_abc123",
}
```

The TypeScript half builds and its output is committed, but it is **not currently installable**:
the `#ts` branch the package expects has never been published, and there are no TypeScript
consumers in the suite. See [docs/development.md](docs/development.md#distribution).

```ts
import type { FacileEvent, FacileAgentSession } from "@facile/enveloppe";
import { FACILE_EVENT_VERSION } from "@facile/enveloppe";

const event: FacileEvent<FacileAgentSession> = {
  version: FACILE_EVENT_VERSION,
  app: "Jardin",
  object: "agent_session",
  action: "created",
  facile_id: "fac_abc123",
  payload: session,
  timestamp: new Date().toISOString(),
  idempotency_key: "jardin_agent_session_created_fac_abc123",
};
```

## Structure

```
go/    enveloppe.go is the whole Go contract
ts/    src/index.ts is the whole TypeScript contract; dist/ is committed output
docs/  architecture, configuration, development, the full type reference
```

## Documentation

| Doc | What's in it |
|---|---|
| [Architecture](docs/architecture.md) | Where the envelope sits, ownership per object, the identity model |
| [Configuration](docs/configuration.md) | Compiler settings, version pinning, JSON serialization rules |
| [Development](docs/development.md) | Build, distribution, and how to change the contract safely |
| [API](docs/api.md) | Every event type and every field, in both languages |

---

Part of the [Facile Suite](https://facile.studio) — self-hosted tools for creative studios
and freelancers. One login, zero cloud dependency.
