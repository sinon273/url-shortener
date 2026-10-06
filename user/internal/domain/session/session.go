package session_domain

import "time"

type Session struct {
	ID     string
	UserID string
	//срок сессии
	ExpiresAt time.Time
	//отозвана ли сессия вручную? при logout=true
	IsRevoked bool
}

func NewSession(ID string, userID string, expiresAt time.Time) Session {
	return Session{
		ID:        ID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}
}

func (s Session) IsAlive() bool {
	if s.IsRevoked {
		return false
	}
	return s.ExpiresAt.After(time.Now())
}
