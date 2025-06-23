package database

import (
	"idm/inner/common"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// ConnectDb получить конфиг и подключиться с ним к базе данных
func ConnectDb(envFile string) (*sqlx.DB, error) {
	cfg, err := common.GetConfig(envFile)
	if err != nil {
		return nil, err
	}

	return ConnectDbWithCfg(cfg)
}

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом
func ConnectDbWithCfg(cfg common.Config) (*sqlx.DB, error) {

	db, err := sqlx.Connect(cfg.DbDriverName, cfg.Dsn)
	if err != nil {
		return nil, err
	}
	db.Close()

	db = sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db, nil
}

func CheckDbConnection(cfg common.Config) bool {
	db, err := sqlx.Connect(cfg.DbDriverName, cfg.Dsn)
	if err != nil {
		return false
	}
	db.Close()
	return true
}
