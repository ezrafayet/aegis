package domain

const (
	RoleUser = "user"
)

type Role struct {
	UserID string `json:"user_id" gorm:"not null;uniqueIndex:idx_user_role"`
	Role   string `json:"role" gorm:"not null;uniqueIndex:idx_user_role"`
}

func NewRole(userID, role string) Role {
	return Role{
		UserID: userID,
		Role:   role,
	}
}

// This is handled by User repository
