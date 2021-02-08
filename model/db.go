package model

import (
	"Blog/util"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var once sync.Once
var err error
var logger = util.GetLogger()

func GetDb() *gorm.DB {
	once.Do(func() {
		initDb()
	})
	return db
}

func initDb() {
	logger.Info("Init db")
	db, err = gorm.Open(util.Db, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		util.DbUser, util.DbPassword, util.DbHost, util.DbPort, util.DbName,
	))
	if err != nil {
		logger.Fatal("Fail to init DB:", err.Error())
	}
	//defer db.Close()
	// 最小链接数（永远保持）
	db.DB().SetMaxIdleConns(10)
	// 最大链接数（链接上限）
	db.DB().SetMaxOpenConns(100)
	// 最大可复用时间 (链接寿命) (不能大于gin框架本身规定的数据库链接时间)
	db.DB().SetConnMaxLifetime(10 * time.Second)
	// 不要复数表名
	db.SingularTable(true)
	db.AutoMigrate(&User{}, &Post{}, &Category{})
}
