package service

import (
	"ginana/internal/service/user"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db   *gorm.DB
	User user.IUser
}

func New(db *gorm.DB, u user.IUser) (s *Service, err error) {
	s = &Service{
		db:   db,
		User: u,
	}
	return
}

func (s *Service) Close() {
	_ = s.db.Close()
}
