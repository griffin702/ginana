package paladin

import (
	"context"
	"errors"
)

var (
	// DefaultClient default client.
	DefaultClient Client
)

// Init init config client.
// If confPath is set, it inits file client by default
// Otherwise we could pass args to init remote client
// args[0]: driver name, string type
func Init(local bool, args ...interface{}) (err error) {
	if len(args) == 0 {
		panic(errors.New("lack of config center args"))
	}
	if local { // use local configs
		configPath, ok := args[0].(string)
		if !ok {
			panic(errors.New("invalid config path"))
		}
		DefaultClient, err = NewFile(configPath)
		return
	}
	var driver Driver
	driverName, ok := args[0].(string)
	if !ok {
		panic(errors.New("invalid remote config center args"))
	}
	driver, err = GetDriver(driverName)
	if err != nil {
		return
	}
	DefaultClient, err = driver.New()
	return
}

// Watch watch on a key. The configuration implements the setter interface, which is invoked when the configuration changes.
func Watch(key string, s Setter) error {
	v := DefaultClient.Get(key)
	str, err := v.Raw()
	if err != nil {
		return err
	}
	if err := s.Set(str); err != nil {
		return err
	}
	go func() {
		for event := range WatchEvent(context.Background(), key) {
			s.Set(event.Value)
		}
	}()
	return nil
}

// WatchEvent watch on multi keys. Events are returned when the configuration changes.
func WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	return DefaultClient.WatchEvent(ctx, keys...)
}

// Get return value by key.
func Get(key string) *Value {
	return DefaultClient.Get(key)
}

// GetAll return all config map.
func GetAll() *Map {
	return DefaultClient.GetAll()
}

// Keys return values key.
func Keys() []string {
	return DefaultClient.GetAll().Keys()
}

// Close close watcher.
func Close() error {
	return DefaultClient.Close()
}
