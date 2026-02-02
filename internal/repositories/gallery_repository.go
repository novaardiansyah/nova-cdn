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

func (r *GalleryRepository) FindAllPaginated(page, limit int, subject_id string, subject_type string, size string) ([]models.Gallery, error) {
	var galleries []models.Gallery
	offset := (page - 1) * limit
	query := r.db.Offset(offset).Limit(limit)

	if subject_id != "" {
		query = query.Where("subject_id = ?", subject_id)
	}

	if subject_type != "" {
		query = query.Where("subject_type = ?", subject_type)
	}

	if size != "" {
		query = query.Where("size = ?", size)
	}

	err := query.Find(&galleries).Error
	return galleries, err
}

func (r *GalleryRepository) Count(subject_id string, subject_type string, size string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Gallery{})

	if subject_id != "" {
		query = query.Where("subject_id = ?", subject_id)
	}

	if subject_type != "" {
		query = query.Where("subject_type = ?", subject_type)
	}

	if size != "" {
		query = query.Where("size = ?", size)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *GalleryRepository) Create(gallery *models.Gallery) error {
	err := r.db.Create(gallery).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GalleryRepository) CreateMany(galleries []*models.Gallery) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, gallery := range galleries {
			if err := tx.Create(gallery).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GalleryRepository) FindByID(id uint64, withDeleted bool) (*models.Gallery, error) {
	var gallery models.Gallery
	if withDeleted {
		err := r.db.Unscoped().First(&gallery, id).Error
		return &gallery, err
	}
	err := r.db.First(&gallery, id).Error
	return &gallery, err
}

func (r *GalleryRepository) Delete(gallery *models.Gallery) error {
	err := r.db.Delete(gallery).Error
	return err
}

func (r *GalleryRepository) ForceDelete(gallery *models.Gallery) error {
	err := r.db.Unscoped().Delete(gallery).Error
	return err
}

func (r *GalleryRepository) Restore(gallery *models.Gallery) error {
	err := r.db.Unscoped().Model(gallery).Update("deleted_at", nil).Error
	return err
}

func (r *GalleryRepository) RestoreByGroupCode(groupCode string, size string) error {
	query := r.db.Unscoped().Model(&models.Gallery{}).Where("group_code = ?", groupCode)

	if size != "" {
		query = query.Where("size = ?", size)
	}

	return query.Update("deleted_at", nil).Error
}

func (r *GalleryRepository) FindByGroupCode(groupCode string, size string) ([]models.Gallery, error) {
	var galleries []models.Gallery
	query := r.db.Where("group_code = ?", groupCode)

	if size != "" {
		query = query.Where("size = ?", size)
	}

	err := query.Find(&galleries).Error

	return galleries, err
}

func (r *GalleryRepository) DeleteByGroupCode(groupCode string, size string) error {
	query := r.db.Where("group_code = ?", groupCode)

	if size != "" {
		query = query.Where("size = ?", size)
	}

	return query.Delete(&models.Gallery{}).Error
}

func (r *GalleryRepository) ForceDeleteByGroupCode(groupCode string, size string) ([]models.Gallery, error) {
	var galleries []models.Gallery
	query := r.db.Unscoped().Where("group_code = ?", groupCode)

	if size != "" {
		query = query.Where("size = ?", size)
	}

	if err := query.Find(&galleries).Error; err != nil {
		return nil, err
	}

	if err := query.Delete(&models.Gallery{}).Error; err != nil {
		return nil, err
	}

	return galleries, nil
}
