// internal/domain/widgets.go
package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/robstave/go-template/internal/domain/types"
)

func (s *Service) GetWidgetByID(widgetID string) (*types.Widget, error) {
	s.logger.Info("Retrieving widget by ID", "widgetID", widgetID)

	widget, err := s.repo.GetWidgetById(widgetID)
	if err != nil {
		s.logger.Error("Error retrieving widget", "error", err)
		return nil, err
	}

	if widget == nil {
		s.logger.Warn("Widget not found", "widgetID", widgetID)
		return nil, errors.New("widget not found")
	}

	return widget, nil
}

func (s *Service) GetAllWidgets() ([]types.Widget, error) {
	s.logger.Info("Retrieving all widgets")

	widgets, err := s.repo.GetAllWidgets()
	if err != nil {
		s.logger.Error("Error retrieving all widgets", "error", err)
		return nil, err
	}

	return widgets, nil
}

func (s *Service) CreateWidget(widget types.Widget) (*types.Widget, error) {
	s.logger.Info("Creating new widget", "name", widget.Name)

	// Generate UUID if not provided
	if widget.ID == "" {
		widget.ID = uuid.New().String()
	}

	if err := s.repo.CreateWidget(widget); err != nil {
		s.logger.Error("Failed to create widget", "error", err)
		return nil, err
	}

	return &widget, nil
}

func (s *Service) UpdateWidget(widget types.Widget) (*types.Widget, error) {
	s.logger.Info("Updating widget", "id", widget.ID)

	// Check if widget exists
	existingWidget, err := s.repo.GetWidgetById(widget.ID)
	if err != nil {
		s.logger.Error("Error checking widget existence", "error", err)
		return nil, err
	}

	if existingWidget == nil {
		s.logger.Warn("Widget not found for update", "id", widget.ID)
		return nil, errors.New("widget not found")
	}

	if err := s.repo.UpdateWidget(widget); err != nil {
		s.logger.Error("Failed to update widget", "error", err)
		return nil, err
	}

	return &widget, nil
}

func (s *Service) DeleteWidget(widgetID string) error {
	s.logger.Info("Deleting widget", "id", widgetID)

	// Check if widget exists
	existingWidget, err := s.repo.GetWidgetById(widgetID)
	if err != nil {
		s.logger.Error("Error checking widget existence", "error", err)
		return err
	}

	if existingWidget == nil {
		s.logger.Warn("Widget not found for deletion", "id", widgetID)
		return errors.New("widget not found")
	}

	if err := s.repo.DeleteWidget(widgetID); err != nil {
		s.logger.Error("Failed to delete widget", "error", err)
		return err
	}

	return nil
}
