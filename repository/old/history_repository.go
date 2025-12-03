package repository_old

import (
	"context"
	"log"
	"moneybkd/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HistoryRepository interface {
	Insert(ctx context.Context, h *models.History) error
	GetByCode(ctx context.Context, code string) ([]*models.History, error)
}

type historyRepository struct {
	col *mongo.Collection
}

func NewHistoryRepository(db *mongo.Database) HistoryRepository {
	return &historyRepository{
		col: db.Collection("history"),
	}
}

func (r *historyRepository) Insert(ctx context.Context, h *models.History) error {
	_, err := r.col.InsertOne(ctx, h)
	return err
}

func (r *historyRepository) GetByCode(ctx context.Context, code string) ([]*models.History, error) {
	log.Println("GET history 4")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"country_code": code}}},
		{{
			Key: "$project",
			Value: bson.M{
				"country_code": 1,
				"country_name": 1,
				"value":        1,
				"day": bson.M{
					"$dateTrunc": bson.M{
						"date": "$created_at",
						"unit": "day",
					},
				},
			},
		}},
		{{
			Key: "$group",
			Value: bson.M{
				"_id": bson.M{
					"country_code": "$country_code",
					"day":          "$day",
				},
				"country_name": bson.M{"$first": "$country_name"},
				"value":        bson.M{"$first": "$value"},
			},
		}},
		{{
			Key: "$project",
			Value: bson.M{
				"_id":          0,
				"country_code": "$_id.country_code",
				"created_at":   "$_id.day",
				"country_name": 1,
				"value":        1,
			},
		}},
		{{Key: "$sort", Value: bson.M{"created_at": 1}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []*models.History

	log.Println("GET history 5")
	for cursor.Next(ctx) {
		var h models.History
		if err := cursor.Decode(&h); err != nil {
			return nil, err
		}
		history = append(history, &h)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	log.Println("GET history 6")

	return history, nil
}
