package main

import (
	"fmt"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := common.GetConfig(".env")
	if err != nil {
		panic(err.Error())
	}

	// создаём подключение к базе данных
	database, err := database.ConnectDbWithCfg(cfg)
	if err != nil {
		panic(err.Error())
	}
	// закрываем соединение с базой данных после выхода из функции main
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}()
	var server = build(database, cfg)
	err = server.App.Listen(":8080")
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

// buil функция, конструирующая наш веб-сервер
func build(database *sqlx.DB, cfg common.Config) *web.Server {
	// создаём веб-сервер
	var server = web.NewServer()
	// создаём репозиторий
	var employeeRepo = employee.NewEmployeeRepository(database)
	var roleRepo = role.NewRoleRepository(database)
	// создаём валидатор
	var vld = validator.NewRequestValidator()
	// создаём сервис
	var employeeService = employee.NewService(employeeRepo, vld)
	var roleService = role.NewService(roleRepo, vld)
	var connectionService = &info.Service{}
	// создаём контроллер
	var employeeController = employee.NewController(server, employeeService)
	var roleController = role.NewController(server, roleService)
	var infoController = info.NewController(server, cfg, connectionService)
	employeeController.RegisterRoutes()
	roleController.RegisterRoutes()
	infoController.RegisterRoutes()

	return server
}
