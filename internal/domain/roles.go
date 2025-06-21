package domain

const (
	RoleUser = "user"
)

type Role struct {
	UserID string `json:"user_id" gorm:"not null;uniqueIndex:idx_user_role"`
	Value  string `json:"role" gorm:"not null;uniqueIndex:idx_user_role"`
	User   User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func NewRole(userID, role string) Role {
	return Role{
		UserID: userID,
		Value:  role,
	}
}

// This is handled by User repository
