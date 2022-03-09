package migrate

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

type Handler struct {
	service Servicer
}

type Servicer interface {
	Up(req *Request, filePath map[string]string) error
	Down(req *Request, filePath map[string]string) error
}

func NewHandler(service Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

var filePath map[string]string
var migrations []Migrations

func (h *Handler) UpDB(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		//return err.JSON(c, errs.NewBadRequest("invalid request", err.Error()))
	}

	initDefault()
	if req.Migrations == nil {
		req.Migrations = migrations
	}

	err := h.service.Up(&req, filePath)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Database up migrated")
}

func (h *Handler) DownDB(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return errors.New("invalid request")
	}

	if len(req.Migrations) <= 0 {
		return errors.New("invalid request - bad request")

	}
	initDefault()
	err := h.service.Down(&req, filePath)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Database down migrated")
}

func initDefault() {
	migrations = []Migrations{
		{
			Type: "data",
		},
		{
			Type: "schema",
		},
	}

	filePath = map[string]string{
		"schema": "file://migration/payments-new",
		"data":   "file://migration/data/" + viper.GetString("migrationenv"),
	}
}
