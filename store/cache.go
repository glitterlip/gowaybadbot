package store

import (
	"encoding/json"
	"github.com/coocood/freecache"
)

var Cache *freecache.Cache

func Put(key string, value interface{}, ttl int) error {
	marshaled, _ := json.Marshal(value)
	err := Cache.Set([]byte(key), marshaled, ttl)
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) {
	Cache.Del([]byte(key))
}
