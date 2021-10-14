package adapters

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"shorty"
)

func DBConnection() *gorm.DB {
	conf := shorty.AppConfig()
	database := conf.Db.DbName
	log.Println("Connecting to adapters")
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to adapters")
	}
	return db
}

func InitialMigrations(database string) {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	db, err := gorm.Open(sqlite.Open(database), config)
	log.Println("migrations running...")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalln("failed to connect adapters")
	}
	err = db.AutoMigrate(&Link{}, &DataPoints{})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalln("failed to run migrations on adapters")
	}
	log.Println("finished running migrations!")
}
