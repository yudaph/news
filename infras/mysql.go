package infras

import (
	"fmt"
	"news/configs"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func MysqlNewClient(cfg configs.Config) (*sqlx.DB, error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DB.MySQL.Username,
		cfg.DB.MySQL.Password,
		cfg.DB.MySQL.Host, cfg.DB.MySQL.Port, cfg.DB.MySQL.Name)
	fmt.Println(ds)
	client, err := sqlx.Open("mysql", ds)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("name", cfg.DB.MySQL.Username).
			Str("host", cfg.DB.MySQL.Host).
			Str("port", cfg.DB.MySQL.Port).
			Str("dbName", cfg.DB.MySQL.Name).
			Msg("Failed connecting to database")
	} else {
		log.
			Info().
			Err(err).
			Str("name", cfg.DB.MySQL.Username).
			Str("host", cfg.DB.MySQL.Host).
			Str("port", cfg.DB.MySQL.Port).
			Str("dbName", cfg.DB.MySQL.Name).
			Msg("Connected to database")
	}

	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 30)
	client.SetConnMaxIdleTime(time.Minute * 1)
	client.SetMaxOpenConns(100)
	client.SetMaxIdleConns(10)

	return client, nil

}
