package role

import (
	"errors"
	"idm/inner/common"
	"idm/inner/web"

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
func (c *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/role"
	c.server.GroupApiV1.Post("/role", c.CreateRole)
	c.server.GroupApiV1.Get("/role", c.GetAllRole)
	c.server.GroupApiV1.Get("/role/id/:id", c.FindRoleById)
	c.server.GroupApiV1.Get("/role/ids/:ids", c.FindRoleByIds)
	c.server.GroupApiV1.Delete("/role/id/:id", c.DeleteRoleById)
	c.server.GroupApiV1.Delete("/role/ids/:ids", c.DeleteRoleByIds)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/role"
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
	var req RequestById
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var foundResponse, err = contr.roleervice.FindById(req.Id)
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
	var req RequestByIds
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var foundResponses, err = contr.roleervice.FindByIds(req.Ids)
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
	var req RequestById
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var err = contr.roleervice.DeleteById(req.Id)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx, 204); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete role")
		return
	}
}

func (contr *Controller) DeleteRoleByIds(ctx *fiber.Ctx) {
	var req RequestByIds
	if err := ctx.QueryParser(req); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}

	var err = contr.roleervice.DeleteByIds(req.Ids)
	if err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if err = common.ResponseWithoutData(ctx, 204); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning result delete roles")
		return
	}
}
