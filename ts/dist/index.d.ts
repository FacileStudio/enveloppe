export type FacileApp = "Sablier" | "Opus" | "Ardoise" | "Plume" | "Glouton" | "Vision";
export type FacileAction = "created" | "updated" | "deleted";
export type FacileObjectType = "project" | "task" | "user" | "time_entry" | "invoice" | "document";
export interface FacileEvent<T = unknown> {
    app: FacileApp;
    object: FacileObjectType;
    action: FacileAction;
    facile_id: string;
    payload: T;
    timestamp: string;
    idempotency_key: string;
}
export interface FacileProject {
    facile_id: string;
    name: string;
    description: string | null;
    icon: string | null;
}
export interface FacileTask {
    facile_id: string;
    project_facile_id: string;
    name: string;
    status: string;
    actor_email?: string;
}
export interface FacileUser {
    facile_id: string;
    email: string;
    name: string;
}
export interface FacileTimeEntry {
    facile_id: string;
    task_facile_id: string;
    user_facile_id: string;
    started_at: string;
    stopped_at: string | null;
}
export interface FacileInvoice {
    facile_id: string;
    client_name: string;
    amount: number;
    currency: string;
    status: string;
}
export interface FacileDocument {
    facile_id: string;
    title: string;
    status: string;
    signer_email: string;
}
