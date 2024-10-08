package infra

import (
	"fmt"
	"os"

	"DellTradingApi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ! Type is not used?
type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func OpenDbConnection() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	sslMode := os.Getenv("DB_SSLMODE")

	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s ", host, user, password, dbName, sslMode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("db err: ", err)
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(10)

	DB = db
	return db
}

// creates/updates tables to database
func Migrate(database gorm.DB) {
	// database.Migrator().DropTable(&models.UserEntity{}, &models.WatchlistEntity{}, &models.PortfolioEntity{})
	// database.Migrator().DropTable(&models.PortfolioEntity{})
	database.AutoMigrate(&models.UserEntity{}, &models.WatchlistEntity{}, &models.PortfolioEntity{})
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := DB.DB()
	sqlDB.Close()
	return err
}

func GetDB() *gorm.DB {
	return DB
}
