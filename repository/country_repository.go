package repository

import (
	"context"
	"encoding/json"
	"log"
	"moneybkd/models"

	supabase "github.com/supabase-community/supabase-go"
)

type CountryRepository interface {
	FindByCode(ctx context.Context, code string) (*models.Country, error)
	Insert(ctx context.Context, c *models.Country) error
	Update(ctx context.Context, c *models.Country) error
	GetAll(ctx context.Context) ([]*models.Country, error)
}

type countryRepo struct {
	client *supabase.Client
}

func NewCountryRepository(c *supabase.Client) CountryRepository {
	return &countryRepo{client: c}
}

func (r *countryRepo) GetAll(ctx context.Context) ([]*models.Country, error) {
	log.Println("Get all countries")
	data, _, err := r.client.From("countries").Select("*", "", false).Execute()
	log.Println("Data")
	log.Println(data)

	var rows []models.Country
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}

	result := make([]*models.Country, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, err
}

func (r *countryRepo) FindByCode(ctx context.Context, code string) (*models.Country, error) {

	data, _, err := r.client.From("countries").
		Select("*", "exact", false).
		Eq("code", code).
		Limit(1, "").
		Execute()

	if err != nil {
		return nil, err
	}

	log.Println("RAW: ", string(data))
	var rows []models.Country
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil // no encontrado
	}

	return &rows[0], nil
}

func (r *countryRepo) Insert(ctx context.Context, c *models.Country) error {
	_, _, err := r.client.From("countries").Insert(c, false, "", "", "").Execute()

	return err
}

func (r *countryRepo) Update(ctx context.Context, c *models.Country) error {
	_, _, err := r.client.From("countries").
		Update(map[string]any{
			"value": c.Value,
		}, "", "").Eq("code", c.Code).
		Execute()

	return err
}
