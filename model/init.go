package model

import (
	"fmt"

	"github.com/lexkong/log"
	"github.com/spf13/viper"
	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dbUsername     = "db.username"
	dbPassword     = "db.password"
	dbAddr         = "db.addr"
	dbName         = "db.name"
	dockerUsername = "docker.username"
	dockerPassword = "docker.password"
	dockerAddr     = "docker.addr"
	dockerName     = "docker.name"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

var DB *Database

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	// db.DB().SetMaxOpenConns(20000)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetMaxIdleConns(0)
}

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		// "Asia/Shanghai"),
		"Local",
	)

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	}

	// set for db connection
	setupDB(db)

	return db
}

// InitSelfDB used for cli
func InitSelfDB() *gorm.DB {
	return openDB(
		viper.GetString(dbUsername),
		viper.GetString(dbPassword),
		viper.GetString(dbAddr),
		viper.GetString(dbName),
	)
}

func InitDockerDB() *gorm.DB {
	return openDB(
		viper.GetString(dockerUsername),
		viper.GetString(dockerPassword),
		viper.GetString(dockerAddr),
		viper.GetString(dockerName),
	)
}

func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

func GetDockerDB() *gorm.DB {
	return InitDockerDB()
}

func (db *Database) Init() {
	DB = &Database{
		Self:   GetSelfDB(),
		Docker: GetDockerDB(),
	}
}

func (db *Database) Close() {
	_ = DB.Self.Close()
	_ = DB.Docker.Close()
}
