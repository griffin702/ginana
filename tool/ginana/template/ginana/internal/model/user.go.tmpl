package model

import (
	"time"
)

type User struct {
	ID          int64      `json:"id" gorm:"primary_key;comment:'用户ID'"`
	CreatedAt   time.Time  `json:"created_at" gorm:"comment:'创建时间'"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"comment:'更新时间'"`
	DeletedAt   *time.Time `json:"-" sql:"index" gorm:"comment:'删除时间戳'"`
	Username    string     `json:"username" gorm:"type:VARCHAR(191);unique;not null;comment:'用户账号'"`
	Password    string     `json:"-" gorm:"type:VARCHAR(255);not null;comment:'用户密码'"`
	Nickname    string     `json:"nickname" gorm:"type:VARCHAR(100);unique;not null;comment:'用户昵称'"`
	Avatar      string     `json:"avatar" gorm:"type:VARCHAR(255);not null;default:'/static/upload/default/user-default-60x60.png';comment:'用户头像'"`
	IsAuth      bool       `json:"is_auth" gorm:"comment:'认证(0-正常,1-未认证)'"`
	LastLogin   time.Time  `json:"last_login" gorm:"type:DATETIME;not null;comment:'最后登录时间'"`
	LastLoginIP string     `json:"last_login_ip" gorm:"type:VARCHAR(30);not null;comment:'最后登录IP'"`
	CountLogin  int64      `json:"count_login" gorm:"comment:'登录次数'"`
	Roles       []*Role    `json:"roles" gorm:"many2many:user_roles"`
}

type UserRoles struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}
