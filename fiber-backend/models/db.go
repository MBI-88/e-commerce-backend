package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	fileLog     *os.File
	loggerError = log.New(fileLog, "DataBase Error => ", log.Ldate)
)

// Make connections
func SetConnection(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	return db
}

// Create error logs by every endpoint
func createLogs(path string) {
	fileLog, err := os.Create(path) 
	if err != nil {
		panic(err)
	}
	loggerError.SetOutput(fileLog)
}

// Migrate database
func Migrate(url string) {
	DB = SetConnection(url)
	if err := DB.AutoMigrate(&Users{},
		&Publishings{}, &Comments{}, &Trolley{},
		&Views{}, &ConfirmationCode{}); err != nil {
		panic(err)
	}
	fmt.Println("[+] Migration done!")
}

// Do connections
func DialDb(url string,logpath string) {
	DB = SetConnection(url)
	go createLogs(logpath)
}
