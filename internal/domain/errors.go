package domain

import "errors"

var (
	ErrGeneric = errors.New("an error occured")
)

var (
	ErrNoUser               = errors.New("no_user")
	ErrUserBlocked          = errors.New("user_blocked")
	ErrUserNotApproved      = errors.New("user_not_approved")
	ErrUserDeleted          = errors.New("user_deleted")
)

var (
	ErrNoRefreshToken       = errors.New("no_refresh_token")
	ErrTooManyRefreshTokens = errors.New("too_many_refresh_tokens")
	ErrInvalidAccessToken   = errors.New("invalid_access_token")
	ErrAccessTokenExpired   = errors.New("access_token_expired")
	ErrRefreshTokenExpired  = errors.New("refresh_token_expired")
)

var (
	ErrWrongAuthMethod      = errors.New("wrong_auth_method")
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
)
