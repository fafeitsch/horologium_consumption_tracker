package orm

import (
	"errors"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
)

type SeriesService interface {
	Save(series *domain.Series) error
	Delete(id uint) error
	QueryById(id uint) (*domain.Series, error)
	QueryAll() ([]domain.Series, error)
}

func NewSeriesService(db *gorm.DB) SeriesService {
	return &seriesServiceImpl{db: db}
}

type seriesServiceImpl struct {
	db *gorm.DB
}

func (s *seriesServiceImpl) Save(series *domain.Series) error {
	entity := seriesEntity{
		Name: series.Name,
	}
	err := s.db.Save(&entity).Error
	series.Id = entity.Id
	return err
}

func (s *seriesServiceImpl) Delete(id uint) error {
	if id == 0 {
		return errors.New("cannot delete entity with id 0")
	}
	entity := seriesEntity{
		Id: id,
	}
	return s.db.Delete(&entity).Error
}

func (s *seriesServiceImpl) QueryById(id uint) (*domain.Series, error) {
	entity := &seriesEntity{}
	err := s.db.First(&entity, id).Error
	series := entity.toDomainSeries()
	return &series, err
}

func (s *seriesServiceImpl) QueryAll() ([]domain.Series, error) {
	resultSet := make([]seriesEntity, 0, 0)
	err := s.db.Find(&resultSet).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.Series, 0, len(resultSet))
	for _, row := range resultSet {
		result = append(result, domain.Series{Id: row.Id, Name: row.Name})
	}
	return result, nil
}
