	package db

	import (
		"fmt"
		"log"

		"github.com/sabarish-manoharan/emp-management/models"
		"github.com/spf13/viper"
		"gorm.io/driver/postgres"
		"gorm.io/gorm"
	)

	var DB *gorm.DB // setting db globally

	func ConnectDB() {
		
	viper.AutomaticEnv() ;
		db_url := getDBurl()

		db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{}) // open up the db

		if err != nil {
			log.Fatalf("Failed to connect DB : %v", err)
		}
		DB = db

		db.AutoMigrate(&models.Employee{})
		db.AutoMigrate(&models.User{})
		fmt.Println("Connected to DB")
	}

	func getDBurl() string {
		return viper.Get("DB_URL").(string) // getting db url from the env
	}
