package database

import (
	"database/sql"
	"fmt"
	glog "github.com/griffin702/ginana/library/log"
	xtime "github.com/griffin702/ginana/library/time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

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
	LogLevel    logger.LogLevel
	Active      int
	Idle        int
	IdleTimeout xtime.Duration
}

type ormLog struct{}

func (l ormLog) Printf(format string, args ...interface{}) {
	if glog.GetLogger() != nil {
		glog.Debugf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *SQLConfig, noDBName ...bool) (db *gorm.DB, err error) {
	if len(noDBName) > 0 && noDBName[0] {
		c.DbName = ""
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.DbUser, c.DbPwd, c.DbHost, c.DbPort, c.DbName, c.Params)
	dialector := mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})
	if c.LogLevel < 1 || c.LogLevel > 4 { // LogLevel 默认为Info
		c.LogLevel = 2
	}
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(ormLog{}, logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      c.LogLevel,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.DbPrev,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Printf("db dsn(%s) error: %v", dsn, err)
		return
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		log.Printf("sqlDB error: %v", err)
		return
	}
	sqlDB.SetMaxIdleConns(c.Idle)
	sqlDB.SetMaxOpenConns(c.Active)
	sqlDB.SetConnMaxLifetime(time.Duration(time.Duration(c.IdleTimeout).Seconds()))
	if c.Debug {
		db = db.Debug()
	}
	return
}
