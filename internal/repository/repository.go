package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	SpendingRepository *SpendingRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		SpendingRepository: NewSpendingRepository(db),
	}
}
