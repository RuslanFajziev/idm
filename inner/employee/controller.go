package employee

import (
	"errors"
	"idm/inner/common"
	"idm/inner/web"

	"github.com/gofiber/fiber"
)

type Controller struct {
	server          *web.Server
	employeeService Srv
}

// интерфейс сервиса employee.Service
type Srv interface {
	FindById(id int64) (Response, error)
	SaveTx(req Request) (id int64, err error)
	FindByIds(ids []int64) ([]Response, error)
	GetAll() ([]Response, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewController(server *web.Server, employeeService Srv) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

// функция для регистрации маршрутов
func (contr *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/employees"
	contr.server.GroupApiV1.Post("/employees", contr.CreateEmployee)
	contr.server.GroupApiV1.Get("/employees", contr.GetAllEmployee)
	contr.server.GroupApiV1.Get("/employees/id/:id", contr.FindEmployeeById)
	contr.server.GroupApiV1.Get("/employees/ids/:ids", contr.FindEmployeeByIds)
	contr.server.GroupApiV1.Delete("/employees/id/:id", contr.DeleteEmployeeById)
	contr.server.GroupApiV1.Delete("/employees/ids/:ids", contr.DeleteEmployeeByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (contr *Controller) CreateEmployee(ctx *fiber.Ctx) {

	// анмаршалим JSON body запроса в структуру Request
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	// вызываем метод SaveTx сервиса employee.Service
	var newId, err = contr.employeeService.SaveTx(req)
	if err != nil {
		switch {

		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
		return
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newId); err != nil {

		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return
	}
}

func (contr *Controller) FindEmployeeById(ctx *fiber.Ctx) {
	var req RequestById
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var foundResponse, err = contr.employeeService.FindById(req.Id)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponse); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employee")
		return
	}
}

func (contr *Controller) FindEmployeeByIds(ctx *fiber.Ctx) {
	var req RequestByIds
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var foundResponses, err = contr.employeeService.FindByIds(req.Ids)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employees")
		return
	}
}

func (contr *Controller) GetAllEmployee(ctx *fiber.Ctx) {
	var foundResponses, err = contr.employeeService.GetAll()
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning all employees")
		return
	}
}

func (contr *Controller) DeleteEmployeeById(ctx *fiber.Ctx) {
	var req RequestById
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var err = contr.employeeService.DeleteById(req.Id)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx, 204); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employee")
		return
	}
}

func (contr *Controller) DeleteEmployeeByIds(ctx *fiber.Ctx) {
	var req RequestByIds
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var err = contr.employeeService.DeleteByIds(req.Ids)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx, 204); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employees")
		return
	}
}
