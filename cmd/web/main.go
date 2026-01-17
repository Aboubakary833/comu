package main

import (
	"comu/config"
	"comu/internal/modules/auth"
	"comu/internal/modules/post"
	"comu/internal/modules/users"
	"comu/internal/shared/logger"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	logger := logger.NewLogger()
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
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(db, config, usersModule.GetPublicApi(), logger)
	postModule := post.NewModule(db, authModule.GetPublicApi(), logger)

	e := echo.New()
	e.Use(
		middleware.Secure(),
		middleware.Recover(),
		middleware.RequestLogger(),
		middleware.RemoveTrailingSlash(),
	)

	// Register modules routes
	authModule.RegisterRoutes(e)
	postModule.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(config.AppAddr))
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
