package enveloppe

type App string

const (
	AppSablier  App = "sablier"
	AppOpus     App = "opus"
	AppCharles  App = "charles"
	AppPlume    App = "plume"
	AppGlouton  App = "glouton"
	AppVision   App = "vision"
)

type Action string

const (
	ActionCreated Action = "created"
	ActionUpdated Action = "updated"
	ActionDeleted Action = "deleted"
)

type ObjectType string

const (
	ObjectProject   ObjectType = "project"
	ObjectTask      ObjectType = "task"
	ObjectUser      ObjectType = "user"
	ObjectTimeEntry ObjectType = "time_entry"
	ObjectInvoice   ObjectType = "invoice"
	ObjectDocument  ObjectType = "document"
)

type Event[T any] struct {
	App            App        `json:"app"`
	Object         ObjectType `json:"object"`
	Action         Action     `json:"action"`
	FacileID       string     `json:"facile_id"`
	Payload        T          `json:"payload"`
	Timestamp      string     `json:"timestamp"`
	IdempotencyKey string     `json:"idempotency_key"`
}
