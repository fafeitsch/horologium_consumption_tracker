package orm

import (
	"errors"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
)

type SeriesService interface {
	Create(series *domain.Series) error
	Delete(id uint) error
	QueryAll() ([]domain.Series, error)
}

func NewSeriesService(db *gorm.DB) SeriesService {
	return &SeriesServiceImpl{db: db}
}

type SeriesServiceImpl struct {
	db *gorm.DB
}

func (s *SeriesServiceImpl) Create(series *domain.Series) error {
	entity := seriesEntity{
		Name: series.Name,
	}
	err := s.db.Save(&entity).Error
	series.Id = entity.ID
	return err
}

func (s *SeriesServiceImpl) Delete(id uint) error {
	if id == 0 {
		return errors.New("cannot delete entity with id 0")
	}
	entity := seriesEntity{
		Model: gorm.Model{ID: id},
	}
	return s.db.Delete(&entity).Error
}

func (s *SeriesServiceImpl) QueryAll() ([]domain.Series, error) {
	resultSet := make([]seriesEntity, 0, 0)
	err := s.db.Find(&resultSet).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.Series, 0, len(resultSet))
	for _, row := range resultSet {
		result = append(result, domain.Series{Id: row.ID, Name: row.Name})
	}
	return result, nil
}
