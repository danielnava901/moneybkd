package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"moneybkd/models"
	"moneybkd/repository"
	"net/http"
	"time"
)

type CurrencyService interface {
	GetCurrency(ctx context.Context, code string) (*models.Country, error)
	UpdateDB() error
	InitCountriesInsert() error
	GetCountries(ctx context.Context) ([]*models.Country, error)
	GetHistoryByCode(ctx context.Context, code string, filter string) ([]*models.History, error)
}

type APIResponse struct {
	Data map[string]float64 `json:"data"`
}

type CurrenciesAPIResponse struct {
	Data map[string]struct {
		Symbol        string `json:"symbol"`
		Name          string `json:"name"`
		SymbolNative  string `json:"symbol_native"`
		DecimalDigits int    `json:"decimal_digits"`
		Rounding      int    `json:"rounding"`
		Code          string `json:"code"`
		NamePlural    string `json:"name_plural"`
	} `json:"data"`
}

type currencyService struct {
	countryRepo repository.CountryRepository
	historyRepo repository.HistoryRepository
	apiKey      string
}

func NewCurrencyService(c repository.CountryRepository, h repository.HistoryRepository, apiKey string) CurrencyService {
	return &currencyService{
		countryRepo: c,
		historyRepo: h,
		apiKey:      apiKey,
	}
}

func (s *currencyService) GetCurrency(ctx context.Context, code string) (*models.Country, error) {
	if code == "" {
		return nil, errors.New("missing country code")
	}

	c, err := s.countryRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errors.New("country not found")
	}

	/*
		h := models.History{
			CountryName: c.Name,
			CountryCode: c.Code,
			Value:       c.Value,
			CreatedAt:   time.Now(),
		}
		s.historyRepo.Insert(ctx, &h)
	*/
	return c, nil
}

func (s *currencyService) GetHistoryByCode(ctx context.Context, code string, filter string) ([]*models.History, error) {
	if code == "" {
		return nil, errors.New("missing country code")
	}

	history, err := s.historyRepo.GetByCode(ctx, code, filter)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (s *currencyService) UpdateDB() error {
	log.Println("corriendo cron...")
	log.Println(time.Now())
	url := "https://api.freecurrencyapi.com/v1/latest?apikey=" + s.apiKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var data APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	ctx := context.Background()

	for code, value := range data.Data {
		log.Println("code: ", code, " value: ", value)
		country, err := s.countryRepo.FindByCode(ctx, code)

		if err != nil {
			continue
		}

		history := models.History{
			CountryName: country.Name,
			CountryCode: code,
			Value:       value,
			CreatedAt:   time.Now(),
		}

		if err := s.historyRepo.Insert(ctx, &history); err != nil {
			return err
		}

		log.Println("Actualizando country ", country.Name, country.Code)

		err = s.countryRepo.Update(ctx, &models.Country{
			Code:  country.Code,
			Value: value,
		})

		if err != nil {
			return err
		}
	}

	log.Println("fin cron...")
	return nil
}

func (s *currencyService) InitCountriesInsert() error {
	url := "https://api.freecurrencyapi.com/v1/currencies?apikey=" + s.apiKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response CurrenciesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	log.Println("response init insert....")
	log.Println(response)

	ctx := context.Background()

	for code, item := range response.Data {
		exists, err := s.countryRepo.FindByCode(ctx, code)
		if err != nil {
			return err
		}

		if exists != nil {
			continue
		}

		country := models.Country{
			Name:      item.Name,
			Code:      item.Code,
			Symbol:    item.Symbol,
			Value:     0,
			UpdatedAt: time.Now(),
		}

		if err := s.countryRepo.Insert(ctx, &country); err != nil {
			return err
		}
	}
	return nil
}

func (s *currencyService) GetCountries(ctx context.Context) ([]*models.Country, error) {
	c, err := s.countryRepo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return c, nil
}
