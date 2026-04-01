package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// ClassRepository defines interface for class CRUD.
type ClassRepository interface {
	Create(class *models.Class) error
	GetByID(id string) (*models.Class, error)
	List() ([]models.Class, error)
	Update(class *models.Class) error
	Delete(id string) error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) GetByID(id string) (*models.Class, error) {
	var class models.Class
	if err := r.db.Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) List() ([]models.Class, error) {
	var classes []models.Class
	if err := r.db.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id string) error {
	return r.db.Delete(&models.Class{}, "id = ?", id).Error
}
