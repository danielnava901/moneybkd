package repository

import (
	"context"
	"encoding/json"
	"moneybkd/models"

	supabase "github.com/supabase-community/supabase-go"
)

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

	var rows []models.History
	if err := json.Unmarshal([]byte(resp), &rows); err != nil {
		return nil, err
	}

	result := make([]*models.History, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}
