package api

import (
	"log"
	"moneybkd/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("INIT HANDLER...............")
	var echoServer = server.New()
	echoServer.ServeHTTP(w, r)
	log.Println("END INIT HANDLER...............")
}
