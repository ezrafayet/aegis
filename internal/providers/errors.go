package providers

import "errors"

var (
	ErrorNoUser          = errors.New("no_user")
	ErrorUserBlocked     = errors.New("user_blocked")
	ErrorUserDeleted     = errors.New("user_deleted")
	ErrorWrongAuthMethod = errors.New("wrong_auth_method")
)

var (
	AuthMethodNotEnabled = errors.New("auth_method_not_enabled")
)