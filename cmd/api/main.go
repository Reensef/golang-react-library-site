package main

import (
	"log"

	"github.com/Reensef/golang-react-boolib/internal/db"
	"github.com/Reensef/golang-react-boolib/internal/env"
	"github.com/Reensef/golang-react-boolib/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := appConfig{
		appAddr: ":" + env.GetString("APP_PORT"),
		sqlDBConfig: sqlDBConfig{
			host:         env.GetString("SQL_DB_HOST"),
			port:         env.GetString("SQL_DB_PORT"),
			name:         env.GetString("SQL_DB_NAME"),
			user:         env.GetString("SQL_DB_USER"),
			password:     env.GetString("SQL_DB_PASSWORD"),
			maxOpenConns: env.GetInt("SQL_DB_MAX_OPEN_CONNS"),
			maxIdleConns: env.GetInt("SQL_DB_MAX_IDLE_CONNS"),
			maxIdleTime:  env.GetString("SQL_DB_MAX_IDLE_TIME"),
		},
		blobDBConfig: blobDBConfig{
			addr: env.GetString("BLOB_DB_HOST") + ":" + env.GetString("BLOB_DB_PORT"),
			id:   env.GetString("BLOB_DB_ID"),
			key:  env.GetString("BLOB_DB_KEY"),
		},
		env: env.GetString("ENV"),
	}

	sqlDB, err := db.NewSql(
		"postgres://"+cfg.sqlDBConfig.user+":"+cfg.sqlDBConfig.password+"@"+cfg.sqlDBConfig.host+":"+cfg.sqlDBConfig.port+"/"+cfg.sqlDBConfig.name+"?sslmode=disable",
		cfg.sqlDBConfig.maxOpenConns,
		cfg.sqlDBConfig.maxIdleConns,
		cfg.sqlDBConfig.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer sqlDB.Close()

	blobDB, err := db.NewBlob(
		cfg.blobDBConfig.addr,
		cfg.blobDBConfig.id,
		cfg.blobDBConfig.key,
	)
	if err != nil {
		log.Panic(err)
	}

	store := store.NewStorage(sqlDB, blobDB)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
