package main

import (
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"
	"gitea.nefixestrada.com/nefix/urlshortener/pkg/handler"
)

func main() {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		log.Fatalf("error opening the DB: %v", err)
	}
	defer func() {
		if err = boltDB.Close(); err != nil {
			log.Fatalf("error closing the DB connection: %v", err)
		}
	}()

	db := &db.DB{
		DB: boltDB,
	}

	if err := db.Initialize(); err != nil {
		log.Fatalf("error initializing the DB: %v", err)
	}

	log.Println("Starting to listen at port :3000")
	if err := http.ListenAndServe(":3000", handler.Default(db)); err != nil {
		log.Fatalf("error listening: %v", err)
	}
}
