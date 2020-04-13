package service

import (
	"context"
	"fmt"
	"ginana/internal/config"
	"ginana/library/cache/memcache"
	"ginana/library/database"
	"ginana/library/tools"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Close()
	SetEnforcer(ef *casbin.SyncedEnforcer) (err error)
	GetAllRoles(ctx context.Context) (roles []database.CasbinRole, err error)
	GetAllUsers(ctx context.Context) (roles []database.CasbinUser, err error)
}

func New(cfg *config.Config, db *gorm.DB, mc memcache.Memcache) (s Service, err error) {
	s = &service{
		cfg:  cfg,
		db:   db,
		mc:   mc,
		tool: tools.New(),
	}
	return
}

type service struct {
	cfg  *config.Config
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	mc   memcache.Memcache
	tool *tools.Tool
}

func (s *service) Close() {
	_ = s.db.Close()
}

// Close close the resource.
func (s *service) SetEnforcer(ef *casbin.SyncedEnforcer) (err error) {
	if !s.cfg.Casbin.Enable {
		return
	}
	if s.tool.PtrIsNil(ef) {
		return fmt.Errorf("enforcer is nil")
	}
	s.ef = ef
	return
}
