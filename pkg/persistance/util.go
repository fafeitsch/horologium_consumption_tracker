package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

func createInMemoryDb() (*gorm.DB, error) {
	originalTableNameHandler := gorm.DefaultTableNameHandler
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		defaultTableName = strings.Replace(defaultTableName, "_entities", "", -1)
		return originalTableNameHandler(db, defaultTableName)
	}
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	db = db.Set("gorm:auto_preload", true)
	err = db.AutoMigrate(&seriesEntity{}, &pricingPlanEntity{}, &meterReadingEntity{}).
		Exec("PRAGMA foreign_keys = ON").
		Error
	return db, err
}
