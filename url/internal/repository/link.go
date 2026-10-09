package url_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	url_domain "github.com/sinon273/url-shortener/url/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "short_code", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})
	if err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, link url_domain.Link) error {
	model := linkDomainToModel(link)
	_, err := r.coll.InsertOne(ctx, model)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return url_domain.ErrAlreadyExists
		}
		return fmt.Errorf("insert link: %w", err)
	}
	return nil
}

func (r *Repository) GetByShortCode(ctx context.Context, shortCode string) (url_domain.Link, error) {
	filter := bson.M{"short_code": shortCode}
	var model LinkModel
	err := r.coll.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return url_domain.Link{}, url_domain.ErrNotFound
		}
		return url_domain.Link{}, fmt.Errorf("find link by short_code: %w", err)
	}
	return linkModelToDomain(model), nil
}

func (r *Repository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	filter := bson.M{"short_code": shortCode}
	res, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete link by short_code: %w", err)
	}
	if res.DeletedCount == 0 {
		return url_domain.ErrNotFound
	}
	return nil
}

func (r *Repository) ListByUser(ctx context.Context, userID string, limit int, cursor string) ([]url_domain.Link, string, error) {
	filter := bson.M{"user_id": userID}

	if cursor != "" {
		oid, err := bson.ObjectIDFromHex(cursor)
		if err != nil {
			return nil, "", url_domain.ErrInvalidArgument
		}
		filter["_id"] = bson.M{"$gt": oid}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(int64(limit + 1))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, "", fmt.Errorf("find link by user_id: %w", err)
	}
	defer cur.Close(ctx)

	var models []LinkModel
	if err = cur.All(ctx, &models); err != nil {
		return nil, "", fmt.Errorf("decode links: %w", err)
	}

	var nextCursor string
	if len(models) > limit {
		nextCursor = models[limit-1].ID.Hex()
		models = models[:limit]
	}

	links := make([]url_domain.Link, len(models))
	for i, m := range models {
		links[i] = linkModelToDomain(m)
	}

	return links, nextCursor, nil
}

func (r *Repository) CountActiveByUser(ctx context.Context, userID string) (int, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
		"$or": []bson.M{
			{"expires_at": bson.M{"$exists": false}},
			{"expires_at": bson.M{"$gt": time.Now()}},
		},
	}

	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count active by links: %w", err)
	}
	return int(count), nil
}
