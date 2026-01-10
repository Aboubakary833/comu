package main

import (
	"comu/config"
	"comu/internal/shared"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	logger := shared.NewLogger()
	config, err := config.NewConfig()

	if err != nil {
		logger.Error.Fatalln(err)
	}

	db, err := openDB(config.DBDriver, config.DBSource)

	if err != nil {
		logger.Error.Fatal(err)
	}
	defer db.Close()

	// Initialize modules and inject db and logging dependencies

	router := chi.NewRouter().With(
		middleware.Logger,
		middleware.Recoverer,
		middleware.CleanPath,
		middleware.RedirectSlashes,
	)
	// Register modules routes

	logger.Info.Printf("Server listening on %s\n", config.AppAddr)

	if err := http.ListenAndServe(config.AppAddr, router); err != nil {
		logger.Error.Fatalln(err.Error())
	}
}

func openDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
