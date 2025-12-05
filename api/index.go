package api

import (
	"moneybkd/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var echoServer = server.New()
	echoServer.ServeHTTP(w, r)
}
