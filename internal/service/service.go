package service

import (
	"ginana/internal/service/i_user"
	"ginana/library/cache/memcache"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	mc   memcache.Memcache
	User i_user.IUser
}

func New(
	db *gorm.DB,
	ef *casbin.SyncedEnforcer,
	mc memcache.Memcache,
	u i_user.IUser,
) (s *Service, err error) {
	if err = u.SetEnforcer(ef); err != nil {
		return
	}
	s = &Service{
		db:   db,
		ef:   ef,
		mc:   mc,
		User: u,
	}
	return
}

func (s *Service) Close() {
	_ = s.db.Close()
}
