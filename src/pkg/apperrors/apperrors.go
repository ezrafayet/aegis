package apperrors

import "errors"

var (
	ErrGeneric = errors.New("an_error_occured")
)

var (
	ErrNoUser               = errors.New("no_user")
	ErrUserBlocked          = errors.New("user_blocked")
	ErrUserDeleted          = errors.New("user_deleted")
	ErrEarlyAdoptersOnly    = errors.New("early_adopters_only")
	ErrNameAlreadyExists    = errors.New("name_already_exists")
	ErrNoName               = errors.New("no_name")
	ErrNoEmail              = errors.New("no_email")
	ErrNoRefreshToken       = errors.New("no_refresh_token")
	ErrTooManyRefreshTokens = errors.New("too_many_refresh_tokens")
	ErrAccessTokenInvalid   = errors.New("token_invalid")
	ErrAccessTokenExpired   = errors.New("token_expired")
	ErrRefreshTokenExpired  = errors.New("token_expired")
	ErrWrongAuthMethod      = errors.New("wrong_auth_method")
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
	ErrInvalidState         = errors.New("invalid_state")
)
