package main

import (
	"log"

	"github.com/Reensef/golang-react-boolib/internal/db"
	"github.com/Reensef/golang-react-boolib/internal/env"
	"github.com/Reensef/golang-react-boolib/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := appConfig{
		appAddr: env.GetString("ADDR"),
		sqlDBConfig: sqlDBConfig{
			addr:         env.GetString("SQL_DB_HOST"),
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
		cfg.sqlDBConfig.addr,
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

	store := store.NewStorage(sqlDB, &blobDB)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
