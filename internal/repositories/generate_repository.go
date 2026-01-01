package repositories

import (
	"nova-cdn/internal/models"

	"gorm.io/gorm"
)

type GenerateRepository struct {
	db *gorm.DB
}

func NewGenerateRepository(db *gorm.DB) *GenerateRepository {
	return &GenerateRepository{db: db}
}

func (r *GenerateRepository) FindByAlias(alias string) (*models.Generate, error) {
	var generate models.Generate
	err := r.db.Unscoped().Where("alias = ?", alias).First(&generate).Error
	if err != nil {
		return nil, err
	}
	return &generate, nil
}

func (r *GenerateRepository) Update(generate *models.Generate) error {
	return r.db.Save(generate).Error
}
