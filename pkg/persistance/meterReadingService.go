package orm

import (
	"errors"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
	"time"
)

type MeterReadingService interface {
	Save(*domain.MeterReading) error
	Delete(uint) error
	QueryForSeries(uint) ([]domain.MeterReading, error)
	QueryOpenInterval(uint, time.Time, time.Time) ([]domain.MeterReading, error)
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

//QueryOpenInterval returns all meter readings between start end AND for the specified series.
//Additionally, it also returns the nearest meter reading left of the interval, as wellas the
//nearest meter reading right of the interval, if available.
func (m *MeterReadingServiceImpl) QueryOpenInterval(seriesId uint, start time.Time, end time.Time) ([]domain.MeterReading, error) {
	resultSet := make([]domain.MeterReading, 0, 0)
	first := meterReadingEntity{}
	err := m.db.Where("DATE(date) < ? AND series_id = ?", start, seriesId).Order("date desc").Limit(1).FirstOrInit(&first).Error
	if err != nil {
		return nil, fmt.Errorf("could not query meter reading left of interval [%v, %v] for seriesId %d: %v", start, end, seriesId, err)
	}
	if first != (meterReadingEntity{}) {
		resultSet = append(resultSet, first.toDomainMeterReadingEntity())
	}
	interval := make([]meterReadingEntity, 0, 0)
	err = m.db.Where("DATE(date) >= ? AND DATE(date) <= ? AND series_id = ?", start, end, seriesId).Order("date asc").Find(&interval).Error
	if err != nil {
		return nil, fmt.Errorf("could not query meter readings in interval [%v, %v] for seriesId %d: %v", start, end, seriesId, err)
	}
	for _, e := range interval {
		resultSet = append(resultSet, e.toDomainMeterReadingEntity())
	}
	last := meterReadingEntity{}
	err = m.db.Where("DATE(date) > ? AND series_id = ?", end, seriesId).Order("date asc").Limit(1).FirstOrInit(&last).Error
	if err != nil {
		return nil, fmt.Errorf("could not query meter reading right of interval [%v, %v] for seriesId %d: %v", start, end, seriesId, err)
	}
	if last != (meterReadingEntity{}) {
		resultSet = append(resultSet, last.toDomainMeterReadingEntity())
	}
	return resultSet, nil
}
