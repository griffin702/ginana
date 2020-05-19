package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	glog "github.com/griffin702/ginana/library/log"
	xtime "github.com/griffin702/ginana/library/time"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
	"time"
)

// Config mysql config.
type SQLConfig struct {
	Driver      string
	DbUser      string
	DbPwd       string
	DbName      string
	DbHost      string
	DbPort      string
	Params      string
	DbPrev      string
	Debug       bool
	Active      int
	Idle        int
	IdleTimeout xtime.Duration
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	glog.Infof(strings.Repeat("%v ", len(v)), v...)
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *SQLConfig) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.DbUser, c.DbPwd, c.DbHost, c.DbPort, c.DbName, c.Params)
	db, err = gorm.Open(c.Driver, dsn)
	if err != nil {
		log.Printf("db dsn(%s) error: %v", dsn, err)
		return
	}
	if c.DbPrev != "" {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return c.DbPrev + defaultTableName
		}
		db.SingularTable(true)
	}
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) / time.Second)
	if glog.GetLogger() != nil {
		db.SetLogger(ormLog{})
	}
	// 创建和更新时间钩子
	//db.Callback().Create().Replace("gorm:update_time_stamp",updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp",updateTimeStampForUpdateCallback)
	return
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				_ = createdAtField.Set(nowTime)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				_ = updatedAtField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("UpdatedAt", time.Now().Unix())
	}
}
