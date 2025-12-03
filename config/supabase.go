// config/supabase.go
package config

import (
	"log"
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

var Supabase *supabase.Client

func ConnectSupabase() {
	log.Println("SUPA URL")
	log.Println(os.Getenv("SUPABASE_URL"))

	Supabase, _ = supabase.NewClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		nil,
	)
}
