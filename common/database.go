package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"mouban/model"
	"net/url"
	"time"
)

var Db *gorm.DB

func InitDatabase() {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")

	tryCreateDB(username, password, host, port, database)
	getConnection(username, password, host, port, database, charset, loc)
	migrateTables()
}

func tryCreateDB(username string, password string, host string, port string, database string) {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Infoln("database close failed")
		}
	}(db)

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ;", database))
	if err != nil {
		panic(err)
	}
}

func getConnection(username string, password string, host string, port string, database string, charset string, loc string) {

	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))

	dbLogger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             500 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		logrus.Infoln("Open database failed", err)
		panic("Open database failed " + err.Error())
	}
	Db = db
	logrus.Infoln("mysql connect success")
}

func migrateTables() {

	err := Db.AutoMigrate(
		&model.Access{},
		&model.Book{},
		&model.Comment{},
		&model.Game{},
		&model.Movie{},
		&model.Song{},
		&model.Rating{},
		&model.Schedule{},
		&model.User{},
		&model.Storage{},
	)
	if err != nil {
		panic("init database failed " + err.Error())
	}

}
