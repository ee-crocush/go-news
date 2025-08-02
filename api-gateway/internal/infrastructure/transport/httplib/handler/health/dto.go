package health

// HealthResponse описывает структуру ответа на /health
type HealthResponse struct {
	Status  string `json:"status" example:"OK"`
	Service string `json:"service" example:"api-gateway"`
	Version string `json:"version" example:"1.0.0"`
	Uptime  string `json:"uptime" example:"1h23m45s"`
}
