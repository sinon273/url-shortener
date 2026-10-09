package url_repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(coll *mongo.Collection) *Repository {
	return &Repository{coll: coll}
}
