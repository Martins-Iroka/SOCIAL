package main

import (
	"log"

	"github.com/Martins-Iroka/social/internal/db"
	"github.com/Martins-Iroka/social/internal/env"
	"github.com/Martins-Iroka/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	store := store.NewPostgresStorage(conn)
	db.Seed(store)
}
