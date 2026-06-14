package domain

// HealthStats holds database connection pool statistics returned by the health endpoint.
type HealthStats struct {
	Status            string `json:"status"`
	Message           string `json:"message,omitempty"`
	Error             string `json:"error,omitempty"`
	OpenConnections   int    `json:"open_connections,omitempty"`
	InUse             int    `json:"in_use,omitempty"`
	Idle              int    `json:"idle,omitempty"`
	WaitCount         int64  `json:"wait_count,omitempty"`
	WaitDuration      string `json:"wait_duration,omitempty"`
	MaxIdleClosed     int64  `json:"max_idle_closed,omitempty"`
	MaxLifetimeClosed int64  `json:"max_lifetime_closed,omitempty"`
}
