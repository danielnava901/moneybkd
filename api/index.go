package api

import (
	"log"
	"moneybkd/server"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("INIT HANDLER HOST...............")
	log.Println(r.Host)
	log.Println("METHOD...............")
	log.Println(r.Method)
	log.Println("URL...............")
	log.Println(r.URL)
	log.Println("URL PATH...............")
	log.Println(r.URL.Path)
	var echoServer = server.New()
	echoServer.ServeHTTP(w, r)
	log.Println("END INIT HANDLER...............")
}
