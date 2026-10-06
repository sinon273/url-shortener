package user_auth

import "time"

type Auth struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(secret string, accessTTL time.Duration, refreshTTL time.Duration) *Auth {
	return &Auth{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (a *Auth) RefreshTTL() time.Duration {
	return a.refreshTTL
}

func (a *Auth) AccessTTL() time.Duration {
	return a.accessTTL
}
