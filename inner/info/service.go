package info

import (
	"idm/inner/common"
	"idm/inner/database"

	_ "github.com/lib/pq"
)

func CheckDbConnection(cfg common.Config) bool {
	return database.CheckDbConnection(cfg)
}
