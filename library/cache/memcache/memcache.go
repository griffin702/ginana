package memcache

import (
	"encoding/json"
	mc "github.com/bradfitz/gomemcache/memcache"
	xtime "github.com/griffin702/ginana/library/time"
	"reflect"
	"time"
)

// Config memcache config.
type Config struct {
	Addr        string
	IdleConns   int
	Timeout     xtime.Duration
	CacheExpire xtime.Duration
}

func New(cfg *Config) Memcache {
	client := mc.New(cfg.Addr)
	client.MaxIdleConns = cfg.IdleConns
	client.Timeout = time.Duration(cfg.Timeout)
	return &memcache{
		client:      client,
		cacheExpire: int32(time.Duration(cfg.CacheExpire) / time.Second),
	}
}

type Memcache interface {
	Get(key string, obj interface{}) (err error)
	Set(key string, v interface{}) (err error)
	Add(key string, v interface{}) (err error)
	Replace(key string, v interface{}) (err error)
	GetMulti(keys []string) (items map[string]*mc.Item, err error)
	Touch(key string, seconds int32) (err error)
	Increment(key string, delta uint64) (newValue uint64, err error)
	Decrement(key string, delta uint64) (newValue uint64, err error)
	DeleteAll() (err error)
	Delete(key string) (err error)
	FlushAll() (err error)
	CompareAndSwap(item *mc.Item) (err error)
}

// Memcache memcache client
type memcache struct {
	client      *mc.Client
	cacheExpire int32
}

func (m *memcache) decode(item *mc.Item, v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	data := item.Value
	switch v.(type) {
	case *[]byte:
		d := v.(*[]byte)
		*d = data
	case *string:
		d := v.(*string)
		*d = string(data)
	case interface{}:
		err = json.Unmarshal(data, v)
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
	err = m.client.Set(&mc.Item{Key: key, Value: value, Expiration: m.cacheExpire})
	return
}

func (m *memcache) Add(key string, v interface{}) (err error) {
	value, err := m.encode(v)
	if err != nil {
		return
	}
	err = m.client.Add(&mc.Item{Key: key, Value: value, Expiration: m.cacheExpire})
	return
}

func (m *memcache) Replace(key string, v interface{}) (err error) {
	value, err := m.encode(v)
	if err != nil {
		return
	}
	err = m.client.Replace(&mc.Item{Key: key, Value: value, Expiration: m.cacheExpire})
	return
}

func (m *memcache) GetMulti(keys []string) (items map[string]*mc.Item, err error) {
	items, err = m.client.GetMulti(keys)
	return
}

func (m *memcache) Touch(key string, seconds int32) (err error) {
	err = m.client.Touch(key, seconds)
	return
}

func (m *memcache) Increment(key string, delta uint64) (newValue uint64, err error) {
	newValue, err = m.client.Increment(key, delta)
	return
}

func (m *memcache) Decrement(key string, delta uint64) (newValue uint64, err error) {
	newValue, err = m.client.Decrement(key, delta)
	return
}

func (m *memcache) DeleteAll() (err error) {
	err = m.client.DeleteAll()
	return
}

func (m *memcache) Delete(key string) (err error) {
	err = m.client.Delete(key)
	return
}

func (m *memcache) FlushAll() (err error) {
	err = m.client.FlushAll()
	return
}

func (m *memcache) CompareAndSwap(item *mc.Item) (err error) {
	err = m.client.CompareAndSwap(item)
	return
}
