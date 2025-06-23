package info

import (
	"idm/inner/common"
	"idm/inner/database"

	_ "github.com/lib/pq"
)

type Service struct {
}

func (serv *Service) CheckDbConnection(cfg common.Config) bool {
	return database.CheckDbConnection(cfg)
}
