package main

import (
	"context"
	"fmt"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// Запускаем сервер в отдельной горутине
	go func() {
		err := server.App.Listen(":8080")
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg)

	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	fmt.Println("Graceful shutdown complete.")
}

// Функция "элегантного" завершения работы сервера по сигналу от операционной системы
func gracefulShutdown(server *web.Server, wg *sync.WaitGroup) {
	// Уведомить основную горутину о завершении работы
	defer wg.Done()
	// Создаём контекст, который слушает сигналы прерывания от операционной системы
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	// Слушаем сигнал прерывания от операционной системы
	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Server exiting")
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
