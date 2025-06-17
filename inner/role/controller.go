package role

import (
	"errors"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
	"strings"

	"github.com/gofiber/fiber"
)

type Controller struct {
	server     *web.Server
	roleervice Srv
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

func NewController(server *web.Server, roleervice Srv) *Controller {
	return &Controller{
		server:     server,
		roleervice: roleervice,
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
func (contr *Controller) CreateRole(ctx *fiber.Ctx) {

	// анмаршалим JSON body запроса в структуру Request
	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	// вызываем метод Save сервиса role.Service
	var newId, err = contr.roleervice.Save(req)
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
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created role id")
		return
	}
}

func (contr *Controller) FindRoleById(ctx *fiber.Ctx) {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
		return
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
		return
	}

	foundResponse, err := contr.roleervice.FindById(num)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponse); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found role")
		return
	}
}

func (contr *Controller) FindRoleByIds(ctx *fiber.Ctx) {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
		return
	}

	if req.IDs == "" {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
		return
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	var foundResponses, err = contr.roleervice.FindByIds(ids)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning found roles")
		return
	}
}

func (contr *Controller) GetAllRole(ctx *fiber.Ctx) {
	var foundResponses, err = contr.roleervice.GetAll()
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.OkResponse(ctx, foundResponses); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning all roles")
		return
	}
}

func (contr *Controller) DeleteRoleById(ctx *fiber.Ctx) {
	var idStr string
	if idStr = ctx.Params("id"); idStr == "" {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "error retrieving id")
		return
	}

	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error converted id tot int64")
		return
	}

	err = contr.roleervice.DeleteById(num)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete role")
		return
	}
}

func (contr *Controller) DeleteRoleByIds(ctx *fiber.Ctx) {
	var req struct {
		IDs string `query:"ids"`
	}

	if err := ctx.QueryParser(&req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid query parameters")
		return
	}

	if req.IDs == "" {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "ids parameter is required")
		return
	}

	idStrings := strings.Split(req.IDs, ",")
	var ids []int64

	for _, idStr := range idStrings {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id format: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	var err = contr.roleervice.DeleteByIds(ids)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete roles")
		return
	}
}
