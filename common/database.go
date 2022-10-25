package common

import (
	"database/sql"
	"fmt"
	"log"
	"mouban/model"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func TryCreateDB(username string, password string, host string, port string, database string) {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	db, err := sql.Open("mysql", sqlStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ;", database))
	if err != nil {
		panic(err)
	}
	fmt.Println("Create db success.")
}

func GetConnection(username string, password string, host string, port string, database string, charset string, loc string) *gorm.DB {

	// 字符串拼接
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc),
	)

	// 配置日志输出
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // 缓存日志时间
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(sqlStr), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println("打开数据库失败", err)
		panic("打开数据库失败" + err.Error())
	}
	return db
}

func MigrateTables(db *gorm.DB) {

	err := db.AutoMigrate(
		&model.Access{},
		&model.Book{},
		&model.Comment{},
		&model.Game{},
		&model.Movie{},
		&model.Music{},
		&model.Queue{},
		&model.Rating{},
		&model.User{},
	)
	if err != nil {
		panic("初始化数据库失败" + err.Error())
	}

}

func InitDB() {
	// 从配置文件中获取参数
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")

	TryCreateDB(username, password, host, port, database)
	db := GetConnection(username, password, host, port, database, charset, loc)
	MigrateTables(db)
}

func init() {
}
