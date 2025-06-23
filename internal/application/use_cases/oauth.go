package use_cases

type OAuthGithubService struct {
	Config                 domain.Config
	Provider               domain.OAuthProviderRepository
	UserRepository         domain.UserRepository
	RefreshTokenRepository domain.RefreshTokenRepository
	StateRepository        domain.StateRepository
}

var _ domain.OAuthProviderService = OAuthGithubService{}

func NewOAuthGithubService(c domain.Config, p domain.OAuthProviderRepository, userRepository domain.UserRepository, refreshTokenRepository domain.RefreshTokenRepository, stateRepository domain.StateRepository) OAuthGithubService {
	return OAuthGithubService{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		StateRepository:        stateRepository,
	}
}

// /!\ Can be abused to generate a lot of states - todo: fix
func (s OAuthGithubService) GetAuthURL(redirectUri string) (string, error) {
	state, err := tokengen.Generate("state_", 13)
	if err != nil {
		return "", err
	}
	if err := s.StateRepository.CreateState(domain.NewState(state)); err != nil {
		return "", err
	}
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		s.Config.Auth.Providers.GitHub.ClientID,
		state,
	)
	return redirectURL, nil
}

func (s OAuthGithubService) ExchangeCode(code, state string) (*http.Cookie, *http.Cookie, error) {
	serverState, err := s.StateRepository.GetAndDeleteState(state)
	if err != nil {
		return nil, nil, apperrors.ErrInvalidState
	}
	if serverState.IsExpired() {
		return nil, nil, apperrors.ErrInvalidState
	}
	userInfos, err := s.Provider.GetUserInfos(code, state, "")
	if err != nil {
		return nil, nil, err
	}

	user, err := domain.GetOrCreateUserIfAllowed(s.UserRepository, userInfos, s.Config)
	if err != nil {
		return nil, nil, err
	}

	if user.AuthMethod != "github" {
		return nil, nil, apperrors.ErrWrongAuthMethod
	}

	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := domain.GenerateTokensForUser(user, "device-id", s.Config, s.RefreshTokenRepository)
	if err != nil {
		return nil, nil, err
	}

	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, s.Config)
	return &accessCookie, &refreshCookie, nil
}
