package orm

import (
	"errors"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
)

type MeterReadingService interface {
	Save(*domain.MeterReading) error
	Delete(uint) error
	QueryForSeries(uint) ([]domain.MeterReading, error)
}

func NewMeterReadingService(db *gorm.DB) MeterReadingService {
	return &MeterReadingServiceImpl{db: db}
}

type MeterReadingServiceImpl struct {
	db *gorm.DB
}

func (m *MeterReadingServiceImpl) Save(reading *domain.MeterReading) error {
	entity := toMeterReadingEntity(*reading)
	err := m.db.Save(&entity).Error
	reading.Id = entity.Id
	reading.Series.Id = entity.SeriesID
	return err
}

func (m *MeterReadingServiceImpl) Delete(id uint) error {
	if id == 0 {
		return errors.New("cannot delete entity with id 0")
	}
	entity := meterReadingEntity{
		Id: id,
	}
	return m.db.Delete(&entity).Error
}

func (m *MeterReadingServiceImpl) QueryForSeries(seriesId uint) ([]domain.MeterReading, error) {
	resultSet := make([]meterReadingEntity, 0, 0)
	series := seriesEntity{
		Id: seriesId,
	}
	err := m.db.Model(&series).Related(&resultSet, "SeriesID").Error
	result := make([]domain.MeterReading, 0, len(resultSet))
	for _, res := range resultSet {
		result = append(result, res.toDomainMeterReadingEntity())
	}
	return result, err
}
