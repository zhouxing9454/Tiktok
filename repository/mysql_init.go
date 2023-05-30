package repository

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

func InitMySQL() (err error) {
	dsn := "root:123456@tcp(mysql_a:3306)/byte_dance?charset=utf8mb4&parseTime=True&loc=Local" //mysql_a_1容器里面的byte_dance数据库
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func ModelAutoMigrate() {
	DB.AutoMigrate(&User{}, &Video{}, &Comment{})
}

func Close() error {
	err := DB.Close()
	if err != nil {
		return errors.New("can't close current db")
	}
	return nil
}
