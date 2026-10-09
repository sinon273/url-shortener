package url_repository

import url_domain "github.com/sinon273/url-shortener/url/internal/domain"

func linkModelToDomain(m LinkModel) url_domain.Link {
	return url_domain.Link{
		ShortCode:   m.ShortCode,
		OriginalURL: m.OriginalURL,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		ExpiresAt:   m.ExpiresAt,
		IsActive:    m.IsActive,
	}
}

func linkDomainToModel(l url_domain.Link) LinkModel {
	return LinkModel{
		ShortCode:   l.ShortCode,
		OriginalURL: l.OriginalURL,
		UserID:      l.UserID,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
		ExpiresAt:   l.ExpiresAt,
		IsActive:    l.IsActive,
	}
}
