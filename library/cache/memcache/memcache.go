package memcache

import (
	"encoding/json"
	xtime "ginana/library/time"
	mc "github.com/bradfitz/gomemcache/memcache"
	"reflect"
	"time"
)

// Config memcache config.
type Config struct {
	Addr      string
	IdleConns int
	Timeout   xtime.Duration
}

func New(cfg *Config) Memcache {
	client := mc.New(cfg.Addr)
	client.MaxIdleConns = cfg.IdleConns
	client.Timeout = time.Duration(cfg.Timeout)
	return &memcache{client: client}
}

type Memcache interface {
	Get(key string, obj interface{}) (err error)
	Set(key string, v interface{}) (err error)
}

// Memcache memcache client
type memcache struct {
	client *mc.Client
}

func (m *memcache) decode(item *mc.Item, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	if item != nil {
		err = json.Unmarshal(item.Value, v)
	}
	return
}

func (m *memcache) encode(v interface{}) (value []byte, err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	value, err = json.Marshal(v)
	return
}

func (m *memcache) Get(key string, obj interface{}) (err error) {
	item, err := m.client.Get(key)
	if err != nil {
		return
	}
	err = m.decode(item, obj)
	return
}

func (m *memcache) Set(key string, v interface{}) (err error) {
	value, err := m.encode(v)
	if err != nil {
		return
	}
	err = m.client.Set(&mc.Item{Key: key, Value: value})
	return
}
