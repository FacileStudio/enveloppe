# enveloppe

Canonical event envelope for the Facile Suite. Shared contract between all apps communicating through [Nook](https://github.com/FacileStudio/Nook).

Both TypeScript and Go implementations live in this repo — one source of truth, both languages always in sync.

## Envelope format

```json
{
  "app": "opus",
  "object": "project",
  "action": "created",
  "facile_id": "fac_cuid2xyz",
  "payload": { "name": "Acme Corp", "description": "Website redesign" },
  "timestamp": "2026-05-22T14:30:00Z",
  "idempotency_key": "opus_proj_created_fac_cuid2xyz_1716388200"
}
```

## Canonical objects

| Object | Fields |
|--------|--------|
| `Project` | `facile_id`, `name`, `description` |
| `Task` | `facile_id`, `project_facile_id`, `name`, `status` |
| `User` | `facile_id`, `email`, `name` |
| `TimeEntry` | `facile_id`, `task_facile_id`, `user_facile_id`, `started_at`, `stopped_at` |
| `Invoice` | `facile_id`, `client_name`, `amount`, `currency`, `status` |
| `Document` | `facile_id`, `title`, `status`, `signer_email` |

## TypeScript

```bash
bun add github:FacileStudio/enveloppe#ts
```

```typescript
import type { FacileEvent, FacileProject } from "@facile/enveloppe";
```

## Go

```bash
go get github.com/FacileStudio/enveloppe/go
```

```go
import enveloppe "github.com/FacileStudio/enveloppe/go"

var event enveloppe.Event[enveloppe.Project]
```
