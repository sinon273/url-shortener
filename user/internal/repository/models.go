package user_postgres_repository

import "time"

type UserModel struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SessionModel struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	IsRevoked bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SubscriptionModel struct {
	ID         string
	Plan       string
	LinksLimit int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
