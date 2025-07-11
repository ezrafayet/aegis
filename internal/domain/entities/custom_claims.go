package entities

import (
	"errors"
	"strings"
)

type CustomClaims struct {
	UserID       string   `json:"user_id"`
	EarlyAdopter bool     `json:"early_adopter"`
	Roles        string   `json:"roles"` // coma separated list
	Metadata     string   `json:"metadata"`
}

func NewCustomClaimsFromValues(userID string, earlyAdopter bool, roles []Role, metadata string) (*CustomClaims, error) {
	if userID == "" {
		return nil, errors.New("custom_claims: user_id is required to create custom claims")
	}
	if len(roles) == 0 {
		return nil, errors.New("custom_claims: roles are required to create custom claims")
	}
	if metadata == "" {
		return nil, errors.New("custom_claims: metadata is required to create custom claims")
	}
	rolesString := ""
	for i, role := range roles {
		rolesString += role.Value
		if i < len(roles)-1 {
			rolesString += ","
		}
	}
	return &CustomClaims{
		UserID:       userID,
		EarlyAdopter: earlyAdopter,
		Roles:        rolesString,
		Metadata:     metadata,
	}, nil
}

func NewCusomClaimsFromMap(ccMap map[string]any) (*CustomClaims, error) {
	if ccMap["user_id"] == nil {
		return nil, errors.New("custom_claims: user_id is required to un-map custom claims")
	}
	if ccMap["early_adopter"] == nil {
		return nil, errors.New("custom_claims: early_adopter is required to un-map custom claims")
	}
	if ccMap["roles"] == nil {
		return nil, errors.New("custom_claims: roles are required to un-map custom claims")
	}
	if ccMap["metadata"] == nil {
		return nil, errors.New("custom_claims: metadata is required to un-map custom claims")
	}
	return &CustomClaims{
		UserID:       ccMap["user_id"].(string),
		EarlyAdopter: ccMap["early_adopter"].(bool),
		Roles:        ccMap["roles"].(string),
		Metadata:     ccMap["metadata"].(string),
	}, nil
}

func (cc *CustomClaims) ToMap() map[string]any {
	return map[string]any{
		"user_id":       cc.UserID,
		"early_adopter": cc.EarlyAdopter,
		"roles":         cc.Roles,
		"metadata":      cc.Metadata,
	}
}

func (cc *CustomClaims) GetRoles() []Role {
	roles := strings.Split(cc.Roles, ",")
	rolesFinal := make([]Role, len(roles))
	for i, role := range roles {
		rolesFinal[i] = Role{Value: strings.TrimSpace(role)}
	}
	return rolesFinal
}
