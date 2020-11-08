package main

import (
	"fmt"
	"time"

	"github.com/tkanos/gonfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type grocery struct {
	ID        uint           `json:"id" gorm:"autoIncrement:true"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Name      string         `json:"name" gorm:"primaryKey"`
	Price     float32        `json:"price"`
}

type errorResponse struct {
	Time    time.Time `json:"timestamp"`
	Message string    `json:"message"`
}

type successMessage struct {
	Message string `json:"message"`
}

type configuration struct {
	Host     string
	Port     int
	User     string
	DBname   string
	SSLmode  string
	Password string
}

type groceries []grocery

var db *gorm.DB
var dbErr error

func initialisedb() {
	var initGroceries = groceries{
		{
			Name:  "Protein Bar",
			Price: 2.50,
		},
		{
			Name:  "Yogurt",
			Price: 1.50,
		},
	}
	config := configuration{}
	err := gonfig.GetConf("config/config.development.json", &config)
	if err != nil {
		panic("Couldnt load config file")
	}
	db, dbErr = gorm.Open(postgres.Open(fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=%v password=%v",
		config.Host, config.Port, config.User, config.DBname, config.SSLmode, config.Password)),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if dbErr != nil {

		panic("failed to connect database")

	}
	db.AutoMigrate(&grocery{})
	for _, initGrocery := range initGroceries {
		db.FirstOrCreate(&initGrocery)
	}
	// db.FirstOrInit(&initGroceries)
}
