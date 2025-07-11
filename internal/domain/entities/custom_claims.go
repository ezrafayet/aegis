package entities

import (
	"errors"
	"strings"
)

type CustomClaims struct {
	UserID       string `json:"user_id"`
	EarlyAdopter bool   `json:"early_adopter"`
	Roles        string `json:"roles"` // coma separated list
	Metadata     string `json:"metadata"`
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
	cClaims := CustomClaims{}
	if ccMap["user_id"] == nil {
		cClaims.UserID = ccMap["user_id"].(string)
	}
	if ccMap["early_adopter"] == nil {
		cClaims.EarlyAdopter = ccMap["early_adopter"].(bool)
	}
	if ccMap["roles"] == nil {
		cClaims.Roles = ccMap["roles"].(string)
	}
	if ccMap["metadata"] == nil {
		cClaims.Metadata = ccMap["metadata"].(string)
	}
	return &cClaims, nil
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
