package service

import (
	"errors"
	"fmt"
)

func NewHelperMap() (hm HelperMap, err error) {
	hm = &helperMap{
		ErrorHelper: map[int]string{
			500:  "服务器错误",
			1001: "查询失败",
			1002: "创建失败",
			1003: "更新失败",
			1004: "删除失败",
		},
		CacheKey: map[int]string{
			1: "user",
			2: "role",
		},
	}
	return
}

type HelperMap interface {
	GetError(i int, args ...interface{}) (int, error)
	GetCacheKey(i int, args ...int64) string
}

type helperMap struct {
	ErrorHelper map[int]string
	CacheKey    map[int]string
}

func (hm *helperMap) GetError(i int, args ...interface{}) (int, error) {
	msg := hm.ErrorHelper[i]
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
		if err, ok := arg.(error); ok {
			return i, err
		}
		if str, ok := arg.(string); ok {
			msg = str
		}
	}
	return i, errors.New(msg)
}

func (hm *helperMap) GetCacheKey(i int, args ...int64) string {
	if len(args) > 1 {
		panic("too many arguments")
	}
	key := hm.CacheKey[i]
	if len(args) == 1 {
		key = fmt.Sprintf("%s_%d", key, args[0])
	}
	return key
}
