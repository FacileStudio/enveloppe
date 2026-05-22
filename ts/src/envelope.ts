export type FacileApp =
  | "sablier"
  | "opus"
  | "charles"
  | "plume"
  | "glouton"
  | "vision";

export type FacileAction = "created" | "updated" | "deleted";

export type FacileObjectType =
  | "project"
  | "task"
  | "user"
  | "time_entry"
  | "invoice"
  | "document";

export interface FacileEvent<T = unknown> {
  app: FacileApp;
  object: FacileObjectType;
  action: FacileAction;
  facile_id: string;
  payload: T;
  timestamp: string;
  idempotency_key: string;
}
