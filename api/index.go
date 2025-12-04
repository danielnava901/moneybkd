package api

import (
	"moneybkd/server"
	"net/http"
)

var echo = server.New()

func Handler(w http.ResponseWriter, r *http.Request) {
	echo.ServeHTTP(w, r)
}
