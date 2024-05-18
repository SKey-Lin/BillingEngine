package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"squalux.com/skey/lending/entities"
	"squalux.com/skey/lending/util"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := util.GetEnv("DB_USERNAME") + ":" + util.GetEnv("DB_PASSWORD") + "@tcp(" + util.GetEnv("DB_HOST") + ":" + util.GetEnv("DB_PORT") + ")/" + util.GetEnv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		log.Panicf("Error connect to database: %s", err)
	}

	db.AutoMigrate(&entities.Borrower{}, &entities.Loan{}, &entities.RepaymentSchedule{})

	DB = db
}
