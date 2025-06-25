package employee

import (
	"errors"
	"fmt"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	server          *web.Server
	employeeService Srv
	logger          *common.Logger
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

func NewController(server *web.Server, employeeService Srv, logger *common.Logger) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
		logger:          logger,
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
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// логируем тело запроса
	contr.logger.Debug("create employee: received request", zap.Any("request", req))

	// вызываем метод SaveTx сервиса employee.Service
	var newId, err = contr.employeeService.SaveTx(req)
	if err != nil {
		switch {

		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			contr.logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			contr.logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}

	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newId); err != nil {
		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		contr.logger.Error("error returning created employee id", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}

	return nil
}

func (contr *Controller) FindEmployeeById(ctx *fiber.Ctx) error {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		contr.logger.Error("error retrieving id")
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		contr.logger.Error("error converted id tot int64", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
	}

	foundResponse, err := contr.employeeService.FindById(num)
	if err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponse); err != nil {
		contr.logger.Error("error returning found employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employee")
	}

	return nil
}

func (contr *Controller) FindEmployeeByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		contr.logger.Error("invalid query parameters", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		contr.logger.Error("ids parameter is required", zap.Error(fmt.Errorf("ids parameter is required")))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			contr.logger.Error("invalid id format: "+idStr, zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
		}
		ids = append(ids, id)
	}

	var foundResponses, err = contr.employeeService.FindByIds(ids)
	if err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		contr.logger.Error("error returning found employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found employees")
	}

	return nil
}

func (contr *Controller) GetAllEmployee(ctx *fiber.Ctx) error {
	var foundResponses, err = contr.employeeService.GetAll()
	if err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		contr.logger.Error("error returning all employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning all employees")
	}

	return nil
}

func (contr *Controller) DeleteEmployeeById(ctx *fiber.Ctx) error {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		contr.logger.Error("error retrieving id", zap.Error(fmt.Errorf("error retrieving id")))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		contr.logger.Error("error converted id tot int64", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
	}

	err = contr.employeeService.DeleteById(num)
	if err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		contr.logger.Error("error returning result delete employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employee")
	}

	return nil
}

func (contr *Controller) DeleteEmployeeByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		contr.logger.Error("invalid query parameters", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		contr.logger.Error("ids parameter is required", zap.Error(fmt.Errorf("ids parameter is required")))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			contr.logger.Error("invalid id format: "+idStr, zap.Error(errors.New("invalid id format:"+idStr)))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
		}
		ids = append(ids, id)
	}

	var err = contr.employeeService.DeleteByIds(ids)
	if err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		contr.logger.Error("error returning result delete employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete employees")
	}

	return nil
}
