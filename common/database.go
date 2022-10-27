package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mouban/model/do"
	"net/url"
)

var Db *gorm.DB

func init() {
	// 从配置文件中获取参数
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
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ;", database))
	if err != nil {
		panic(err)
	}
	fmt.Println("Create db success.")
}

func getConnection(username string, password string, host string, port string, database string, charset string, loc string) {

	// 字符串拼接
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))

	db, err := gorm.Open(mysql.Open(sqlStr))

	if err != nil {
		fmt.Println("打开数据库失败", err)
		panic("打开数据库失败" + err.Error())
	}
	Db = db
}

func migrateTables() {

	err := Db.AutoMigrate(
		&do.Access{},
		&do.Book{},
		&do.Comment{},
		&do.Game{},
		&do.Movie{},
		&do.Music{},
		&do.Rating{},
		&do.Record{},
		&do.User{},
	)
	if err != nil {
		panic("初始化数据库失败" + err.Error())
	}

}
