package employee

import (
	"errors"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	contr.server.GroupApiV1.Get("/employees/ids", contr.FindEmployeeByIds)
	contr.server.GroupApiV1.Delete("/employees/id/:id", contr.DeleteEmployeeById)
	contr.server.GroupApiV1.Delete("/employees/ids", contr.DeleteEmployeeByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (contr *Controller) CreateEmployee(ctx *fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру Request
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод SaveTx сервиса employee.Service
	var newId, err = contr.employeeService.SaveTx(req)
	if err != nil {
		switch {

		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newId); err != nil {

		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}

	return nil
}

func (contr *Controller) FindEmployeeById(ctx *fiber.Ctx) error {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
	}

	foundResponse, err := contr.employeeService.FindById(num)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponse); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employee")
	}

	return nil
}

func (contr *Controller) FindEmployeeByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
		}
		ids = append(ids, id)
	}

	var foundResponses, err = contr.employeeService.FindByIds(ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employees")
	}

	return nil
}

func (contr *Controller) GetAllEmployee(ctx *fiber.Ctx) error {
	var foundResponses, err = contr.employeeService.GetAll()
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning all employees")
	}

	return nil
}

func (contr *Controller) DeleteEmployeeById(ctx *fiber.Ctx) error {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
	}

	err = contr.employeeService.DeleteById(num)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employee")
	}

	return nil
}

func (contr *Controller) DeleteEmployeeByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
		}
		ids = append(ids, id)
	}

	var err = contr.employeeService.DeleteByIds(ids)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employees")
	}

	return nil
}
