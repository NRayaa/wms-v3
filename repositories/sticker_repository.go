package repositories

import (
	"wms/models"

	"gorm.io/gorm"
)

// StickerRepository defines interface for sticker CRUD.
type StickerRepository interface {
	Create(sticker *models.Sticker) error
	GetBySlug(slug string) (*models.Sticker, error)
	GetSlugLike(slug string) ([]models.Sticker, error)
	GetByID(id string) (*models.Sticker, error)
	List() ([]models.Sticker, error)
	Update(sticker *models.Sticker) error
	Delete(id string) error
}

// stickerRepository is GORM implementation.
type stickerRepository struct {
	db *gorm.DB
}

// NewStickerRepository constructor.
func NewStickerRepository(db *gorm.DB) StickerRepository {
	return &stickerRepository{db: db}
}

func (r *stickerRepository) Create(sticker *models.Sticker) error {
	return r.db.Create(sticker).Error
}

func (r *stickerRepository) GetBySlug(slug string) (*models.Sticker, error) {
	var sticker models.Sticker
	if err := r.db.Where("slug = ?", slug).First(&sticker).Error; err != nil {
		return nil, err
	}
	return &sticker, nil
}

func (r *stickerRepository) GetSlugLike(slug string) ([]models.Sticker, error) {
	var stickers []models.Sticker
	if err := r.db.Where("slug LIKE ?", slug+"%").Find(&stickers).Error; err != nil {
		return nil, err
	}
	return stickers, nil
}

func (r *stickerRepository) GetByID(id string) (*models.Sticker, error) {
	var sticker models.Sticker
	if err := r.db.Where("id = ?", id).First(&sticker).Error; err != nil {
		return nil, err
	}
	return &sticker, nil
}

func (r *stickerRepository) List() ([]models.Sticker, error) {
	var stickers []models.Sticker
	if err := r.db.Find(&stickers).Error; err != nil {
		return nil, err
	}
	return stickers, nil
}

func (r *stickerRepository) Update(sticker *models.Sticker) error {
	return r.db.Save(sticker).Error
}

func (r *stickerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Sticker{}).Error
}
