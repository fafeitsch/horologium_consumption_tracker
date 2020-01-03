package orm

import (
	"errors"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
)

type PricingPlanService interface {
	Save(*domain.PricingPlan) error
	Delete(uint) error
	QueryAll() ([]domain.PricingPlan, error)
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
