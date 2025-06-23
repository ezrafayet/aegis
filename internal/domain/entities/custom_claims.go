package entities

type CustomClaims struct {
	UserID       string   `json:"user_id"`
	EarlyAdopter bool     `json:"early_adopter"`
	RolesValues  []string `json:"roles"` // coma separated list
	Metadata     string   `json:"metadata"`
}
