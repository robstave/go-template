package domain

import (
	"log/slog"

	"github.com/robstave/go-template/internal/adapters/repositories"
	"github.com/robstave/go-template/internal/domain/types"
)

type Service struct {
	logger *slog.Logger
	repo   repositories.Repository
}

type Domain interface {
	GetWidgetByID(widgetID string) (*types.Widget, error)
	GetAllWidgets() ([]types.Widget, error)
	CreateWidget(widget types.Widget) (*types.Widget, error)
	UpdateWidget(widget types.Widget) (*types.Widget, error)
	DeleteWidget(widgetID string) error
	SeedWidget() error
}

// internal/domain/service.go
// ...
func NewService(logger *slog.Logger, repo repositories.Repository) Domain {
	service := &Service{
		logger: logger,
		repo:   repo,
	}

	// Seed the initial user.   This is called on every startup, but will only create the user if it doesn't already exist
	// To reset the app, just delete the database file  ( assuming you're using the default sqlite3 database )
	if err := service.SeedWidget(); err != nil {
		logger.Error("Failed to seed initial user", "error", err)
	}

	return service
}
