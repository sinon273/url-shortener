package user_service

import (
	"context"
	"time"

	"github.com/sinon273/url-shortener/user/internal/domain/session"
	user_domain "github.com/sinon273/url-shortener/user/internal/domain/user"
)

type Service struct {
	repo  Repository
	cache Cache
	auth  Auth
}

func NewService(repo Repository, cache Cache, auth Auth) *Service {
	return &Service{repo: repo, cache: cache, auth: auth}
}

type Repository interface {
	CreateUserWithSubscription(ctx context.Context, email, passwordHash string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (user_domain.User, error)
	CreateSession(ctx context.Context, userID string, expiresAt time.Time) (string, error)
	GetSession(ctx context.Context, sessionID string) (session_domain.Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
	GetSubscriptionLimit(ctx context.Context, userID string) (int, error)
}

type Cache interface {
	//кладет в кеш какую то сессию
	SetSession(ctx context.Context, sessionID string, userID string, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string) (string, error)
	DelSession(ctx context.Context, sessionID string) error
}

type Auth interface {
	//хэширует пароль для бд
	HashPassword(password string) (string, error)
	//проверяет паароль при входе, точнее хеши
	ComparePassword(hash, password string) error
	//создает jwt 15 мин
	SignAccess(userID, sessionID string) (string, error)
	//jwt токен но долго живущий 7дней
	SignRefresh(userID, sessionID string) (string, error)
	//проверить токен на каждом защищенном запросе
	//проверить refrech token при обновлениее accesstoken(в RefreshToken)
	Parse(token string) (string, string, error)
	RefreshTTL() time.Duration
	AccessTTL() time.Duration
}
