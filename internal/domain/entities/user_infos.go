package entities

// UserInfos is what is returned by the providers (GitHub, Google, etc.)
type UserInfos struct {
	Name   string
	Email  string
	Avatar string
}
