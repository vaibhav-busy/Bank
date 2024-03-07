package database

import (
	"log"
	"os"
	"time"

	pg "github.com/go-pg/pg/v10"
)

var Db *pg.DB

func Connect() *pg.DB {

	opts := &pg.Options{
		User:         "postgres",
		Password:     "5679",
		Addr:         "localhost:5432",
		Database:     "bank",
		DialTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  30 * time.Minute,
		MaxConnAge:   1 * time.Minute,
		PoolSize:     20,
	}

	Db= pg.Connect(opts)

	if Db == nil {
		log.Printf("Error connecting to database")
		os.Exit(100)
	}

	log.Printf("Successfully connected to database")
	
	return Db
}
