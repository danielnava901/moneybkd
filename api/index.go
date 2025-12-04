package api

import (
	"moneybkd/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var echo = server.New()
	echo.ServeHTTP(w, r)
}
