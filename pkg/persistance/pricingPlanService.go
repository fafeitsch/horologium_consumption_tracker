package orm

import (
	"errors"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
	"time"
)

type PricingPlanService interface {
	Save(*domain.PricingPlan) error
	Delete(uint) error
	QueryById(uint) (*domain.PricingPlan, error)
	QueryAll() ([]domain.PricingPlan, error)
	QueryForSeries(uint) ([]domain.PricingPlan, error)
	QueryForTime(uint, time.Time) (*domain.PricingPlan, error)
}

func NewPricingPlanService(db *gorm.DB) PricingPlanService {
	return &PricingPlanServiceImpl{db: db}
}

type PricingPlanServiceImpl struct {
	db *gorm.DB
}

func (p *PricingPlanServiceImpl) Save(plan *domain.PricingPlan) error {
	entity := toPricingPlanEntity(*plan)
	err := p.db.Save(&entity).Error
	plan.Id = entity.Id
	plan.Series.Id = entity.SeriesID
	return err
}

func (p *PricingPlanServiceImpl) Delete(id uint) error {
	if id == 0 {
		return errors.New("cannot delete entity with id 0")
	}
	entity := pricingPlanEntity{
		Id: id,
	}
	return p.db.Delete(&entity).Error
}

func (p *PricingPlanServiceImpl) QueryAll() ([]domain.PricingPlan, error) {
	resultSet := make([]pricingPlanEntity, 0, 0)
	err := p.db.Find(&resultSet).Error
	result := make([]domain.PricingPlan, 0, len(resultSet))
	for _, res := range resultSet {
		result = append(result, res.toDomainPricingPlan())
	}
	return result, err
}

func (p *PricingPlanServiceImpl) QueryForSeries(seriesId uint) ([]domain.PricingPlan, error) {
	resultSet := make([]pricingPlanEntity, 0, 0)
	series := seriesEntity{
		Id: seriesId,
	}
	err := p.db.Model(series).Related(&resultSet, "seriesID").Error
	result := make([]domain.PricingPlan, 0, len(resultSet))
	for _, res := range resultSet {
		result = append(result, res.toDomainPricingPlan())
	}
	return result, err
}

func (p *PricingPlanServiceImpl) QueryForTime(seriesId uint, t time.Time) (*domain.PricingPlan, error) {
	resultSet := make([]pricingPlanEntity, 0, 0)
	err := p.db.Where("DATE(valid_from) <= ? AND DATE(valid_to) >= ? AND series_id = ?", t, t, seriesId).Find(&resultSet).Error
	if err != nil {
		return nil, err
	}
	if len(resultSet) != 1 {
		return nil, fmt.Errorf("there were %d plans for the date %v AND seriesId %d", len(resultSet), t, seriesId)
	}
	result := resultSet[0].toDomainPricingPlan()
	return &result, nil
}

func (p *PricingPlanServiceImpl) QueryById(id uint) (*domain.PricingPlan, error) {
	entity := &pricingPlanEntity{}
	err := p.db.First(&entity, id).Error
	plan := entity.toDomainPricingPlan()
	return &plan, err
}
