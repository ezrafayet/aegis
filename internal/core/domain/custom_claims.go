package domain

type CustomClaims struct {
	UserID       string   `json:"user_id"`
	EarlyAdopter bool     `json:"early_adopter"`
	RolesValues  []string `json:"roles"`
	Metadata     string   `json:"metadata"`
}
