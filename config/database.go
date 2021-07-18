package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (db *gorm.DB, err error) {

	env := loadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.DB_USER, env.DB_PASS, env.DB_HOST, env.DB_NAME)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//Hentikan Program dan Munculkan Error
		log.Fatal(err.Error())
	}

	return db, err
}
