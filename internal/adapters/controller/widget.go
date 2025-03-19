// internal/adapters/controller/widget.go
package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robstave/go-template/internal/domain/types"
)

// CreateWidget handles widget creation
// @Summary Create a new Widget
// @Description Create a new widget
// @Tags Widgets
// @Accept json
// @Produce json
// @Success 201 {object} types.Widget
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /widgets [post]
func (hc *Controller) CreateWidget(c echo.Context) error {
	var widget types.Widget
	if err := c.Bind(&widget); err != nil {
		hc.logger.Error("Failed to bind widget data", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid widget data"})
	}

	// Call the service to create the widget
	createdWidget, err := hc.service.CreateWidget(widget)
	if err != nil {
		hc.logger.Error("Failed to create widget", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create widget"})
	}

	return c.JSON(http.StatusCreated, createdWidget)
}

// GetWidget retrieves a widget by ID
// @Summary Get a widget by ID
// @Description Get a widget by its ID
// @Tags Widgets
// @Produce json
// @Param id path string true "Widget ID"
// @Success 200 {object} types.Widget
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /widgets/{id} [get]
func (hc *Controller) GetWidget(c echo.Context) error {
	id := c.Param("id")

	widget, err := hc.service.GetWidgetByID(id)
	if err != nil {
		hc.logger.Error("Failed to retrieve widget", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Widget not found"})
	}

	return c.JSON(http.StatusOK, widget)
}

// GetAllWidgets retrieves all widgets
// @Summary Get all widgets
// @Description Get all available widgets
// @Tags Widgets
// @Produce json
// @Success 200 {array} types.Widget
// @Failure 500 {object} map[string]string
// @Router /widgets [get]
func (hc *Controller) GetAllWidgets(c echo.Context) error {
	widgets, err := hc.service.GetAllWidgets()
	if err != nil {
		hc.logger.Error("Failed to retrieve widgets", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve widgets"})
	}

	return c.JSON(http.StatusOK, widgets)
}

// UpdateWidget updates an existing widget
// @Summary Update a widget
// @Description Update an existing widget
// @Tags Widgets
// @Accept json
// @Produce json
// @Param id path string true "Widget ID"
// @Success 200 {object} types.Widget
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /widgets/{id} [put]
func (hc *Controller) UpdateWidget(c echo.Context) error {
	id := c.Param("id")

	var widget types.Widget
	if err := c.Bind(&widget); err != nil {
		hc.logger.Error("Failed to bind widget data", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid widget data"})
	}

	// Ensure ID in path matches body
	widget.ID = id

	updatedWidget, err := hc.service.UpdateWidget(widget)
	if err != nil {
		hc.logger.Error("Failed to update widget", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Widget not found or update failed"})
	}

	return c.JSON(http.StatusOK, updatedWidget)
}

// DeleteWidget deletes a widget by ID
// @Summary Delete a widget
// @Description Delete a widget by its ID
// @Tags Widgets
// @Param id path string true "Widget ID"
// @Success 204 {object} nil
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /widgets/{id} [delete]
func (hc *Controller) DeleteWidget(c echo.Context) error {
	id := c.Param("id")

	if err := hc.service.DeleteWidget(id); err != nil {
		hc.logger.Error("Failed to delete widget", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Widget not found or delete failed"})
	}

	return c.NoContent(http.StatusNoContent)
}
