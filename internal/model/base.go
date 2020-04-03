package model

import "time"

type User struct {
	ID          	int64        `json:"id" gorm:"primary_key;comment:'用户ID'"`
	CreatedAt   	int64        `json:"created_at" gorm:"comment:'创建时间戳'"`
	UpdatedAt   	int64        `json:"updated_at" gorm:"comment:'更新时间戳'"`
	DeletedAt   	*time.Time   `json:"-" sql:"index" gorm:"comment:'删除时间戳'"`
	Username    	string       `json:"username" gorm:"type:VARCHAR(191);unique;not null;comment:'用户账号'"`
	Password    	string       `json:"-" gorm:"type:VARCHAR(255);not null;comment:'用户密码'"`
	Nickname    	string       `json:"nickname" gorm:"type:VARCHAR(100);unique;not null;comment:'用户昵称'"`
	IsAuth      	bool         `json:"is_auth" gorm:"comment:'认证(0-正常,1-未认证)'"`
}