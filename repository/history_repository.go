package repository

import (
	"context"
	"encoding/json"
	"log"
	"moneybkd/models"
	"time"

	supabase "github.com/supabase-community/supabase-go"
)

type apiItem struct {
	CountryCode string  `json:"country_code"`
	CreatedAt   string  `json:"created_at"`
	CountryName string  `json:"country_name"`
	Value       float64 `json:"value"`
}

type HistoryRepository interface {
	Insert(ctx context.Context, h *models.History) error
	GetByCode(ctx context.Context, code string) ([]*models.History, error)
}

type historyRepo struct {
	client *supabase.Client
}

func NewHistoryRepository(c *supabase.Client) HistoryRepository {
	return &historyRepo{client: c}
}

func (r *historyRepo) Insert(ctx context.Context, h *models.History) error {
	_, _, err := r.client.From("history").Insert(h, false, "", "", "").Execute()

	return err
}

func (r *historyRepo) GetByCode(ctx context.Context, code string) ([]*models.History, error) {
	// llamada al RPC
	resp := r.client.Rpc("get_history_by_code", "", map[string]any{
		"code_input": code,
	})

	log.Println("Get by code repo:")
	log.Println(resp)
	var items []apiItem
	json.Unmarshal([]byte(resp), &items)

	var histories []*models.History
	for _, it := range items {
		t, _ := time.Parse("2006-01-02", it.CreatedAt)

		histories = append(histories, &models.History{
			CountryCode: it.CountryCode,
			CountryName: it.CountryName,
			Value:       it.Value,
			CreatedAt:   t,
		})
	}

	return histories, nil
}
