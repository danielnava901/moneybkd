package cron

import (
	"fmt"
	"log"
	"moneybkd/config"
	"moneybkd/repository"
	"moneybkd/service"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Your cron job logic here
	config.ConnectSupabase()

	countryRepo := repository.NewCountryRepository(config.Supabase)
	countryHistory := repository.NewHistoryRepository(config.Supabase)

	apiKey := os.Getenv("EXCHANGE_API_KEY")
	svc := service.NewCurrencyService(countryRepo, countryHistory, apiKey)
	if err := svc.UpdateDB(); err != nil {
		log.Println("Cron job error:", err)
		http.Error(w, "cron failed", 500)
		return
	}

	log.Println("Cron job finished.")
	fmt.Fprintln(w, "Cron OK")

}
