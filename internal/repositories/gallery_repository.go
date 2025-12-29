package repositories

import (
	"nova-cdn/internal/models"

	"gorm.io/gorm"
)

type GalleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
	return &GalleryRepository{db: db}
}

func (r *GalleryRepository) FindAllPaginated(page, limit int) ([]models.Gallery, error) {
	var galleries []models.Gallery
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Find(&galleries).Error
	return galleries, err
}

func (r *GalleryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Gallery{}).Count(&count).Error
	return count, err
}

func (r *GalleryRepository) Create(gallery *models.Gallery) error {
	err := r.db.Create(gallery).Error
	if err != nil {
		return err
	}
	r.db.First(gallery, gallery.ID)
	return nil
}

func (r *GalleryRepository) CreateMany(galleries []*models.Gallery) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, gallery := range galleries {
			if err := tx.Create(gallery).Error; err != nil {
				return err
			}
		}
		for _, gallery := range galleries {
			tx.First(gallery, gallery.ID)
		}
		return nil
	})
}
