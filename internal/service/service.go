package service

import (
	"ginana/internal/service/user"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db   *gorm.DB
	ef   *casbin.SyncedEnforcer
	User user.IUser
}

func New(
	db *gorm.DB,
	ef *casbin.SyncedEnforcer,
	u user.IUser,
) (s *Service, err error) {
	if err = u.SetEnforcer(ef); err != nil {
		return
	}
	s = &Service{
		db:   db,
		ef:   ef,
		User: u,
	}
	return
}

func (s *Service) Close() {
	_ = s.db.Close()
}
