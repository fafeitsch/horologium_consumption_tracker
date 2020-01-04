package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"strings"
)

func CreateInMemoryDb() (*gorm.DB, error) {
	configureTableNameHandler()
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	db, err = configureDatabase(db, err)
	return db, err
}

func configureTableNameHandler() {
	originalTableNameHandler := gorm.DefaultTableNameHandler
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		defaultTableName = strings.Replace(defaultTableName, "_entities", "", -1)
		return originalTableNameHandler(db, defaultTableName)
	}
}

func configureDatabase(db *gorm.DB, err error) (*gorm.DB, error) {
	db = db.Set("gorm:auto_preload", true)
	err = db.AutoMigrate(&seriesEntity{}, &pricingPlanEntity{}, &meterReadingEntity{}).
		Exec("PRAGMA foreign_keys = ON").
		Error
	return db, err
}

func ConnectToFileDb(file string) (*gorm.DB, error) {
	configureTableNameHandler()
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	db, err = configureDatabase(db, err)
	return db, err
}
