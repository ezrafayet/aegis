package providers

import "errors"

var (
	ErrNoUser          = errors.New("no_user")
	ErrUserBlocked     = errors.New("user_blocked")
	ErrUserDeleted     = errors.New("user_deleted")
	ErrWrongAuthMethod = errors.New("wrong_auth_method")
	ErrNoRefreshToken  = errors.New("no_refresh_token")
)

var (
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
)
