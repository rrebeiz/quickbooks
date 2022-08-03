package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rrebeiz/quickbooks/internal/data"
	"log"
	"os"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn    string
		pepper string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
}

const version = "1.0.0"

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "the port of the application")
	flag.StringVar(&cfg.env, "environment", "development", "environment, development | production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DSN"), "DB DSN")
	flag.StringVar(&cfg.db.pepper, "db-pepper", "super-secret-pepper", "DB pepper")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(cfg)
	defer db.Close()
	if err != nil {
		errorLog.Println("Failed to connect to the DB.", err)

	}

	models := data.NewModels(db)
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   models,
	}

	err = app.serve()
	if err != nil {
		errorLog.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
