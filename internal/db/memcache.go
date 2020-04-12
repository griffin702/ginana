package db

import (
	"ginana/library/cache/memcache"
	"ginana/library/conf/paladin"
)

func NewMC() (mc memcache.Memcache, err error) {
	var cfg struct {
		Client *memcache.Config
	}
	if err = paladin.Get("memcache.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	mc = memcache.New(cfg.Client)
	return
}
