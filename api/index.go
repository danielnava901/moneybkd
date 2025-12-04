package api

import (
	"log"
	"moneybkd/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("INIT HANDLER...............")
	var echo = server.New()
	echo.ServeHTTP(w, r)
	log.Println("END INIT HANDLER...............")
}
