package apperrors

import "errors"

var (
	ErrGeneric = errors.New("an_error_occured")
)

var (
	ErrAccessTokenInvalid   = errors.New("access_token_invalid")
	ErrAccessTokenExpired   = errors.New("access_token_expired")
	ErrRefreshTokenInvalid  = errors.New("refresh_token_invalid")
	ErrRefreshTokenExpired  = errors.New("refresh_token_expired")
	ErrTooManyRefreshTokens = errors.New("too_many_refresh_tokens")
)

var (
	ErrNoRoles          = errors.New("no_roles")
	ErrUnauthorizedRole = errors.New("unauthorized_role")
)

var (
	ErrNoUser               = errors.New("no_user")
	ErrUserBlocked          = errors.New("user_blocked")
	ErrUserDeleted          = errors.New("user_deleted")
	ErrEarlyAdoptersOnly    = errors.New("early_adopters_only")
	ErrNameAlreadyExists    = errors.New("name_already_exists")
	ErrNoName               = errors.New("no_name")
	ErrNoEmail              = errors.New("no_email")
	ErrWrongAuthMethod      = errors.New("wrong_auth_method")
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
	ErrInvalidState         = errors.New("invalid_state")
)
