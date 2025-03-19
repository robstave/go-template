package controller

import (
	"log/slog"

	"github.com/robstave/go-template/internal/domain"
)

type Controller struct {
	service domain.Domain
	logger  *slog.Logger
}

func NewController(service domain.Domain, logger *slog.Logger) *Controller {
	return &Controller{service: service, logger: logger}
}
