package database

import (
	"fmt"
	"go-bun-chi/entity"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetUpConfigMySql() *gorm.DB {
	err := initToYmal()
	if err != nil {
		panic("Failed connect to database")
	}

	dbUser := viper.GetString("mysql.user")
	dbPass := viper.GetString("mysql.password")
	dbHost := viper.GetString("mysql.host")
	dbName := viper.GetString("mysql.database")
	dbPort := viper.GetString("mysql.port")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			IgnoreRecordNotFoundError: false, // Ignore ErrRecordNotFound error for logger
			// SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel: logger.Info, // Log level
			Colorful: false,       // Disable color
			// Logger:   logger.Default.LogMode(logger.Info),
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})

	fmt.Println(err)

	if err != nil {
		panic("Failed to create connection")
	}

	log.Println("connected")

	if !db.Migrator().HasTable("customers") {
		db.Migrator().CreateTable(&entity.Customer{})

		password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
		cust := &entity.Customer{
			Name:     "test",
			Email:    "test@test.com",
			Password: string(password),
		}

		err := db.Create(cust).Error
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable("orders") {
		db.Migrator().CreateTable(&entity.Orders{})
	}

	if !db.Migrator().HasTable("order_details") {
		db.Migrator().CreateTable(&entity.OrderDetails{})
	}
	return db
}

func CloseDatabaseMysSqlConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}

	dbSQL.Close()

}

func GetSecretKey() (string, error) {
	err := initToYmal()
	if err != nil {
		return "secret key null", err
	}
	return viper.GetString("secret_key_jwt"), nil
}

func GetRefreshSecretKey() (string, error) {
	err := initToYmal()
	if err != nil {
		return "refresh secret key null", nil
	}
	return viper.GetString("secret_key_jwt_refresh"), nil
}

func initToYmal() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("environment")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
