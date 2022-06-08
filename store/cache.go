package store

import (
	"encoding/json"
	"errors"
	"github.com/coocood/freecache"
	"strconv"
)

var Cache *freecache.Cache

func Put(key string, value interface{}, druation ...int) error {
	var ttl int
	var byteValue []byte
	if temp, ok := value.([]byte); ok {
		byteValue = temp
	} else {
		byteValue, _ = json.Marshal(value)
	}
	if len(druation) > 0 {
		ttl = druation[0]
	}
	err := Cache.Set([]byte(key), byteValue, ttl)
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) {
	Cache.Del([]byte(key))
}
func Get(key string, defaultValue ...interface{}) interface{} {
	value, err := Cache.Get([]byte(key))
	if errors.Is(err, freecache.ErrNotFound) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return value
}
func Increment(key string, step ...int) {
	exist, err := Cache.Get([]byte(key))
	var value int
	if len(step) > 0 {
		value = step[0]
	} else {
		value = 1
	}
	if errors.Is(err, freecache.ErrNotFound) {
		Cache.Set([]byte(key), []byte{byte(value)}, 0)
	}
	intVal, _ := strconv.Atoi(string(exist))
	res := intVal + value
	Cache.Set([]byte(key), []byte(strconv.Itoa(res)), 0)
}
