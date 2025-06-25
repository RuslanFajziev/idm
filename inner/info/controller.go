package info

import (
	"idm/inner/common"
	"idm/inner/web"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	server            *web.Server
	cfg               common.Config
	connectionService Srv
}

type Srv interface {
	CheckDbConnection(cfg common.Config) bool
}

func NewController(server *web.Server, cfg common.Config, connectionService Srv) *Controller {
	return &Controller{
		server:            server,
		cfg:               cfg,
		connectionService: connectionService,
	}
}

type InfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) RegisterRoutes() {
	// полный путь будет "/internal/info"
	c.server.GroupInternal.Get("/info", c.GetInfo)
	// полный путь будет "/internal/health"
	c.server.GroupInternal.Get("/health", c.GetHealth)
}

// GetInfo получение информации о приложении
func (c *Controller) GetInfo(ctx *fiber.Ctx) error {
	var err = ctx.Status(fiber.StatusOK).JSON(&InfoResponse{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	})
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
	}

	return nil
}

// GetHealth проверка работоспособности приложения
func (c *Controller) GetHealth(ctx *fiber.Ctx) error {
	result := c.connectionService.CheckDbConnection(c.cfg)
	if result {
		return ctx.Status(fiber.StatusOK).SendString("Healthy")
	} else {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Unhealthy")
	}
}
