package service

import (
	"context"
	"errors"
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
	GetEFRoles(ctx context.Context) (roles []*database.EFRolePolicy, err error)
	GetEFUsers(ctx context.Context) (users []*database.EFUseRole, err error)
}

func New(cfg *config.Config, db *gorm.DB, mc memcache.Memcache, eh *map[int]string) (s Service, err error) {
	s = &service{
		cfg:  cfg,
		db:   db,
		mc:   mc,
		eh:   eh,
		tool: tools.Tools,
	}
	return
}

type service struct {
	cfg  *config.Config
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	mc   memcache.Memcache
	eh   *map[int]string
	tool *tools.Tool
}

func (s *service) Close() {
	_ = s.db.Close()
}

func (s *service) GetError(i int, args ...string) (int, error) {
	if len(args) > 1 {
		panic("too many arguments")
	}
	errHelper := *s.eh
	msg := errHelper[i]
	if len(args) == 1 {
		msg = args[0]
	}
	return i, errors.New(msg)
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
