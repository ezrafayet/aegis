package domain

import "errors"

var (
	ErrNoUser               = errors.New("no_user")
	ErrUserBlocked          = errors.New("user_blocked")
	ErrUserDeleted          = errors.New("user_deleted")
	ErrWrongAuthMethod      = errors.New("wrong_auth_method")
	ErrNoRefreshToken       = errors.New("no_refresh_token")
	ErrTooManyRefreshTokens = errors.New("too_many_refresh_tokens")
)

var (
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
)
