package repository_old

import (
	"context"
	"moneybkd/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CountryRepository interface {
	FindByCode(ctx context.Context, code string) (*models.Country, error)
	Insert(ctx context.Context, country *models.Country) error
	Update(ctx context.Context, country *models.Country) error
	GetAll(ctx context.Context) ([]*models.Country, error)
}

type countryReporsitory struct {
	col *mongo.Collection
}

func NewCountryRepository(db *mongo.Database) CountryRepository {
	return &countryReporsitory{
		col: db.Collection("countries"),
	}
}

func (r *countryReporsitory) GetAll(ctx context.Context) ([]*models.Country, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var countries []*models.Country

	for cursor.Next(ctx) {
		var c models.Country
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		countries = append(countries, &c)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

func (r *countryReporsitory) FindByCode(ctx context.Context, code string) (*models.Country, error) {
	var c models.Country
	err := r.col.FindOne(ctx, bson.M{"code": code}).Decode(&c)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &c, err
}

func (r *countryReporsitory) Insert(ctx context.Context, country *models.Country) error {
	_, err := r.col.InsertOne(ctx, country)
	return err
}

func (r *countryReporsitory) Update(ctx context.Context, country *models.Country) error {
	update := bson.M{
		"$set": bson.M{
			"value": country.Value,
		},
	}

	_, err := r.col.UpdateOne(ctx, bson.M{"code": country.Code}, update)
	return err
}
