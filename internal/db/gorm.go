package db

import (
	"ginana/internal/config"
	"ginana/internal/model"
	"ginana/library/conf/paladin"
	"ginana/library/database"
	"github.com/jinzhu/gorm"
)

func NewDB(cfg *config.Config) (db *gorm.DB, err error) {
	key := "db.toml"
	if err = paladin.Get(key).UnmarshalTOML(cfg); err != nil {
		return
	}
	db, err = database.NewMySQL(cfg.MySQL)
	if err != nil {
		return
	}
	if cfg.MySQL.Debug {
		db = db.Debug()
	}
	initTable(db)
	initTableData(db)
	return
}

func initTable(db *gorm.DB) {
	db.AutoMigrate(
		new(model.User),
	)
}

func initTableData(db *gorm.DB) {

}
