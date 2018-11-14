package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	bolt "go.etcd.io/bbolt"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"
	"gitea.nefixestrada.com/nefix/urlshortener/pkg/handler"
)

type logWriter struct {
	File *os.File
}

func (w *logWriter) Write(b []byte) (int, error) {
	fmt.Print(string(b[:]))

	return w.File.Write(b)
}

func main() {
	// Configure the logging
	f, err := os.OpenFile("urlshortener.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("error closing the log file: %v", err)
		}
	}()

	w := &logWriter{
		File: f,
	}

	log.SetOutput(w)

	// Open the DB and initialize it
	var boltDB *bolt.DB
	boltDB, err = bolt.Open("urlshortener.db", 0600, nil)
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

	// Start the HTTP server
	log.Println("Starting to listen at port :3000")
	if err := http.ListenAndServe(":3000", handler.Default(db)); err != nil {
		log.Fatalf("error listening: %v", err)
	}
}
