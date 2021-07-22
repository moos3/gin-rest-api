package database

import (
	"fmt"
	"os"

	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // configures postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres" // configures postgres driver
	"github.com/joho/godotenv"
	"github.com/moos3/gin-rest-api/database/models"
)

// Initialize initializes the database
func Initialize() (*gorm.DB, error) {
	fmt.Println("Trying to connect to the db")
	//dbConfig := os.Getenv("DB_CONFIG")
	dbConfig := GenerateDBConnectionString()
	db, err := gorm.Open("postgres", dbConfig)
	db.LogMode(true) // logs SQL
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	models.Migrate(db)
	return db, err
}

// GenerateDBConnectionString - returns the proper dsn
func GenerateDBConnectionString() string {
	var (
		dbHost  string
		dbName  string
		dbUser  string
		dbPass  string
		dbPort  string
		dbSsl   string
		dbType  string
		strBase string
		conn    string
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if checkEnv("DB_HOST") {
		dbHost = os.Getenv("DB_HOST")
	}
	if checkEnv("DB_USER") {
		dbUser = os.Getenv("DB_USER")
	}
	if checkEnv("DB_PASS") {
		dbPass = os.Getenv("DB_PASS")
	}
	if checkEnv("DB_NAME") {
		dbName = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_PORT") != "" {
		dbPort = os.Getenv("DB_PORT")
	} else {
		dbPort = "5432"
	}

	if checkEnv(("DB_TYPE")) {
		dbType = os.Getenv("DB_TYPE")
	} else {
		dbType = "postgres"
	}

	if checkEnv(("DB_SSL")) {
		dbSsl = os.Getenv("DB_SSL")
	} else {
		dbSsl = "disable"
	}

	if dbType == "mysql" {
		strBase = `%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True`
		conn = fmt.Sprintf(strBase, dbUser, dbPass, dbHost, dbPort, dbName)
	}

	fmt.Println(dbType)
	if dbType == "postgres" {
		strBase = "host=%s user=%s dbname=%s password=%s port=%s sslmode=%s"
		conn = fmt.Sprintf(strBase, dbHost, dbUser, dbName, dbPass, dbPort, dbSsl)
	}

	return conn
}

// func - checkEnv validate env var that isn't empty
func checkEnv(envVar string) bool {
	if os.Getenv(envVar) == "" {
		fmt.Fprintf(os.Stderr, envVar+" environment variable must be set.\n")
		os.Exit(1)
	}
	return true
}
