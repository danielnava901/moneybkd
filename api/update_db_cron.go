package api

import (
	"fmt"
	"log"
	"moneybkd/config"
	"moneybkd/repository"
	"moneybkd/service"
	"net/http"
	"os"
)

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Your cron job logic here
	fmt.Println("Cron job executed!")
	fmt.Fprintf(w, "Cron job executed successfully!")
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
