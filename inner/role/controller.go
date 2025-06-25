package role

import (
	"errors"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	server     *web.Server
	roleervice Srv
	logger     *common.Logger
}

// интерфейс сервиса employee.Service
type Srv interface {
	FindById(id int64) (Response, error)
	Save(req Request) (id int64, err error)
	FindByIds(ids []int64) ([]Response, error)
	GetAll() ([]Response, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewController(server *web.Server, roleervice Srv, logger *common.Logger) *Controller {
	return &Controller{
		server:     server,
		roleervice: roleervice,
		logger:     logger,
	}
}

// функция для регистрации маршрутов
func (contr *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/roles"
	contr.server.GroupApiV1.Post("/roles", contr.CreateRole)
	contr.server.GroupApiV1.Get("/roles", contr.GetAllRole)
	contr.server.GroupApiV1.Get("/roles/id/:id", contr.FindRoleById)
	contr.server.GroupApiV1.Get("/roles/ids", contr.FindRoleByIds)
	contr.server.GroupApiV1.Delete("/roles/id/:id", contr.DeleteRoleById)
	contr.server.GroupApiV1.Delete("/roles/ids", contr.DeleteRoleByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/roles"
func (contr *Controller) CreateRole(ctx *fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру Request
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		contr.logger.Error(err.Error(), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// логируем тело запроса
	contr.logger.Debug("create role: received request", zap.Any("request", req))

	// вызываем метод Save сервиса role.Service
	var newId, err = contr.roleervice.Save(req)
	if err != nil {
		switch {

		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			contr.logger.Error("create role", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			contr.logger.Error("create role", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newId); err != nil {
		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		contr.logger.Error("error returning created role id", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created role id")
	}

	return nil
}

func (contr *Controller) FindRoleById(ctx *fiber.Ctx) error {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		contr.logger.Error("error returning created role id")
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		contr.logger.Error("error converted id tot int64", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
	}

	foundResponse, err := contr.roleervice.FindById(num)
	if err != nil {
		contr.logger.Error(err.Error())
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponse); err != nil {
		contr.logger.Error("error returning found role", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found role")
	}

	return nil
}

func (contr *Controller) FindRoleByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		contr.logger.Error("invalid query parameters", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		contr.logger.Error("ids parameter is required")
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

	var foundResponses, err = contr.roleervice.FindByIds(ids)
	if err != nil {
		contr.logger.Error(err.Error())
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		contr.logger.Error("error returning found roles", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found roles")
	}

	return nil
}

func (contr *Controller) GetAllRole(ctx *fiber.Ctx) error {
	var foundResponses, err = contr.roleervice.GetAll()
	if err != nil {
		contr.logger.Error(err.Error())
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		contr.logger.Error("error returning all roles", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning all roles")
	}

	return nil
}

func (contr *Controller) DeleteRoleById(ctx *fiber.Ctx) error {
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

	err = contr.roleervice.DeleteById(num)
	if err != nil {
		contr.logger.Error(err.Error())
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())

	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		contr.logger.Error("error returning result delete role", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete role")
	}

	return nil
}

func (contr *Controller) DeleteRoleByIds(ctx *fiber.Ctx) error {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		contr.logger.Error("invalid query parameters", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if req.IDs == "" {
		contr.logger.Error("ids parameter is required")
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

	var err = contr.roleervice.DeleteByIds(ids)
	if err != nil {
		contr.logger.Error(err.Error())
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		contr.logger.Error("error returning result delete roles", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete roles")
	}

	return nil
}
