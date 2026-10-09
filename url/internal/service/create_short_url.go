package url_service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s *Service) CreateShortURL(ctx context.Context, originalURL, userID string, expiresIn *int64) (shortCode, shortURL string, err error) {
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		return "", "", url_domain.ErrInvalidArgument
	}

	limit, err := s.cache.GetLimit(ctx, userID)
	if errors.Is(err, url_domain.ErrNotFound) {
		limit, err = s.userClient.GetLimit(ctx, userID)
		if err != nil {
			return "", "", url_domain.ErrUnavailable
		}
		_ = s.cache.SetLimit(ctx, userID, limit, 5*time.Minute)
	}

	_, err = s.cache.GetCounter(ctx, userID)
	if errors.Is(err, url_domain.ErrNotFound) {
		count, err := s.repo.CountActiveByUser(ctx, userID)
		if err != nil {
			return "", "", url_domain.ErrUnavailable
		}
		_ = s.cache.SetCounter(ctx, userID, count)
	}

	used, err := s.cache.IncrCounter(ctx, userID)
	if err != nil {
		return "", "", err
	}
	if used > limit {
		_ = s.cache.DecrCounter(ctx, userID)
		return "", "", url_domain.ErrResourceExhausted
	}

	var expiresAt *time.Time
	if expiresIn != nil {
		t := time.Now().Add(time.Duration(*expiresIn) * time.Second)
		expiresAt = &t
	}

	for i := 0; i < 3; i++ {
		code, err := generateShortCode()
		if err != nil {
			_ = s.cache.DecrCounter(ctx, userID)
			return "", "", err
		}

		link := url_domain.Link{
			ShortCode:   code,
			OriginalURL: originalURL,
			UserID:      userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ExpiresAt:   expiresAt,
			IsActive:    true,
		}

		err = s.repo.Create(ctx, link)
		if err == nil {
			shortCode = code
			break
		}

		if errors.Is(err, url_domain.ErrAlreadyExists) {
			continue
		}

		_ = s.cache.DecrCounter(ctx, userID)
		return "", "", err
	}

	if shortCode == "" {
		_ = s.cache.DecrCounter(ctx, userID)
		return "", "", fmt.Errorf("failed to generate unique short code")
	}

	shortURL = s.shortURLBase + "/" + shortCode
	return shortCode, shortURL, nil
}

func generateShortCode() (string, error) {
	buf := make([]byte, 7)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}

	code := make([]byte, 7)
	for i, b := range buf {
		code[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(code), nil
}
