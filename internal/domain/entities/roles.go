package entities

const (
	RoleUser          = "user"
	RolePlatformAdmin = "platform_admin"
)

type Role struct {
	UserID string `json:"user_id" gorm:"not null;uniqueIndex:idx_user_role"`
	Value  string `json:"role" gorm:"not null;uniqueIndex:idx_user_role"`
	// relations
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// todo: cascade delete roles on user deletion

func NewRole(userID, role string) Role {
	return Role{
		UserID: userID,
		Value:  role,
	}
}
