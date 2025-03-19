// internal/domain/seeder.go
package domain

import (
	"github.com/google/uuid"
	"github.com/robstave/go-template/internal/domain/types"
)

func (s *Service) SeedWidget() error {
	// Default widget data
	defaultWidgets := []types.Widget{
		{
			ID:    uuid.New().String(),
			Name:  "Sample Widget 1",
			Value: "This is a sample widget to demonstrate functionality",
		},
		{
			ID:    uuid.New().String(),
			Name:  "Sample Widget 2",
			Value: "Another widget example for testing",
		},
	}

	// Check if we already have widgets
	existingWidgets, err := s.repo.GetAllWidgets()
	if err != nil {
		s.logger.Error("Failed to check existing widgets", "error", err)
		return err
	}

	// If we already have widgets, don't seed
	if len(existingWidgets) > 0 {
		s.logger.Info("Database already has widgets, skipping seed")
		return nil
	}

	// Create the sample widgets
	for _, widget := range defaultWidgets {
		if err := s.repo.CreateWidget(widget); err != nil {
			s.logger.Error("Failed to seed widget", "name", widget.Name, "error", err)
			return err
		}
		s.logger.Info("Seeded widget successfully", "id", widget.ID, "name", widget.Name)
	}

	s.logger.Info("Successfully seeded initial widgets")
	return nil
}
