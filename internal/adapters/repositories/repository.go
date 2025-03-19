// internal/adapters/repositories/user.go
package repositories

import (
	"github.com/robstave/go-template/internal/domain/types"
	"gorm.io/gorm"
)

type Repository interface {
	GetWidgetById(id string) (*types.Widget, error)
	GetAllWidgets() ([]types.Widget, error)
	CreateWidget(widget types.Widget) error
	UpdateWidget(widget types.Widget) error
	DeleteWidget(id string) error
}

type RepositorySQLite struct {
	db *gorm.DB
}

func NewRepositorySQLite(db *gorm.DB) Repository {
	return &RepositorySQLite{db: db}
}

func (r *RepositorySQLite) GetWidgetById(id string) (*types.Widget, error) {
	var widget types.Widget
	result := r.db.First(&widget, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &widget, nil
}

func (r *RepositorySQLite) GetAllWidgets() ([]types.Widget, error) {
	var widgets []types.Widget
	result := r.db.Find(&widgets)
	if result.Error != nil {
		return nil, result.Error
	}
	return widgets, nil
}

func (r *RepositorySQLite) CreateWidget(widget types.Widget) error {
	return r.db.Create(&widget).Error
}

func (r *RepositorySQLite) UpdateWidget(widget types.Widget) error {
	return r.db.Save(&widget).Error
}

func (r *RepositorySQLite) DeleteWidget(id string) error {
	return r.db.Delete(&types.Widget{}, "id = ?", id).Error
}
