package user

import (
	"context"
	"fmt"
	"ginana/internal/config"
	"ginana/library/database"
	"ginana/library/tools"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
	"time"
)

// Service service interface
type IUser interface {
	SetEnforcer(ef *casbin.SyncedEnforcer) (err error)
	GetAllRoles(ctx context.Context) (roles []database.CasbinRole, err error)
	GetAllUsers(ctx context.Context) (roles []database.CasbinUser, err error)
}

// Service service
type SUser struct {
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	cfg  *config.Config
	tool *tools.Tool
	//mc              *memcache.Memcache
	cacheExpire int32
	errHelper   *map[int]string
}

// New new a service and return
func New(db *gorm.DB, cfg *config.Config) (s IUser, err error) {
	s = &SUser{
		db:          db,
		cfg:         cfg,
		tool:        tools.New(),
		cacheExpire: int32(time.Duration(cfg.CacheExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (s *SUser) SetEnforcer(ef *casbin.SyncedEnforcer) (err error) {
	if !s.cfg.Casbin.Enable {
		return
	}
	if s.tool.PtrIsNil(ef) {
		return fmt.Errorf("enforcer is nil")
	}
	s.ef = ef
	return
}

func (s *SUser) GetAllRoles(ctx context.Context) (roles []database.CasbinRole, err error) {
	return
}

func (s *SUser) GetAllUsers(ctx context.Context) (roles []database.CasbinUser, err error) {
	return
}
