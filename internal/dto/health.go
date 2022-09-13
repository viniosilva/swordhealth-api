package dto

type HealthStatus string

const (
	HealshStatusUp   HealthStatus = "up"
	HealshStatusDown HealthStatus = "down"
)

type HealthResponse struct {
	Status HealthStatus `json:"status,omitempty"`
}
