package enveloppe

type Project struct {
	FacileID    string  `json:"facile_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Task struct {
	FacileID        string `json:"facile_id"`
	ProjectFacileID string `json:"project_facile_id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
}

type User struct {
	FacileID string `json:"facile_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type TimeEntry struct {
	FacileID     string  `json:"facile_id"`
	TaskFacileID string  `json:"task_facile_id"`
	UserFacileID string  `json:"user_facile_id"`
	StartedAt    string  `json:"started_at"`
	StoppedAt    *string `json:"stopped_at"`
}

type Invoice struct {
	FacileID   string  `json:"facile_id"`
	ClientName string  `json:"client_name"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Status     string  `json:"status"`
}

type Document struct {
	FacileID    string `json:"facile_id"`
	Title       string `json:"title"`
	Status      string `json:"status"`
	SignerEmail string `json:"signer_email"`
}
