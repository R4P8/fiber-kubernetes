package repository

import (
	"example/config"
	"example/entities"
)

type CategoryRepository interface {
	FindAll() ([]entities.Category, error)
	FindByID(id uint) (*entities.Category, error)
	Create(category *entities.Category) error
	Update(category *entities.Category) error
	Delete(id uint) error
}

type categoryRepositoryImpl struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepositoryImpl{}
}

func (r *categoryRepositoryImpl) FindAll() ([]entities.Category, error) {
	var categories []entities.Category
	if err := config.GormDB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepositoryImpl) FindByID(id uint) (*entities.Category, error) {
	var category entities.Category
	if err := config.GormDB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) Create(category *entities.Category) error {
	return config.GormDB.Create(category).Error
}

func (r *categoryRepositoryImpl) Update(category *entities.Category) error {
	return config.GormDB.Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(id uint) error {
	return config.GormDB.Delete(&entities.Category{}, id).Error
}
