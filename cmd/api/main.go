package main

import (
	"log"

	"github.com/Martins-Iroka/social/internal/env"
	"github.com/Martins-Iroka/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewPostgresStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()

	log.Fatal(app.run(mux))
}
