package main

import (
	"log"

	"github.com/Martins-Iroka/social/internal/db"
	"github.com/Martins-Iroka/social/internal/env"
	"github.com/Martins-Iroka/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30), // this set an upper limit of open connection to your connection pool.
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30), // having more takes resources but improves performance
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}
	// important for cleaning resources
	defer db.Close()
	log.Printf("database connection pool established")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()

	log.Fatal(app.run(mux))
}
