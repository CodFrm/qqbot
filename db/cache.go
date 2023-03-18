package db

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/codfrm/cago/database/redis"
)

type Option func(*Options)
type Options struct {
	TTL     time.Duration
	Context context.Context
}

func WithTTL(t time.Duration) Option {
	return func(options *Options) {
		options.TTL = t
	}
}

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, v := range opts {
		v(options)
	}
	return options
}

func GetOrSet(key string, get interface{}, set func() (interface{}, error), opts ...Option) error {
	options := NewOptions(opts...)
	data, err := redis.Default().Get(context.Background(), key).Result()
	if err != nil {
		val, err := set()
		if err != nil {
			return err
		}
		ttl := time.Duration(0)
		if options.TTL > 0 {
			ttl = options.TTL
		} else {
			ttl = time.Hour * 72
		}
		b, err := json.Marshal(val)
		if err != nil {
			return err
		}
		if err := redis.Default().Set(context.Background(), key, b, ttl).Err(); err != nil {
			return err
		}
		copyInterface(get, val)
	} else {
		if err := json.Unmarshal([]byte(data), get); err != nil {
			return err
		}
	}
	return nil
}

type StringCache struct {
	String string
}

type IntCache struct {
	Int int
}

func Get(key string, get interface{}, opts ...Option) error {
	val, err := redis.Ctx(context.Background()).Get(key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), get); err != nil {
		return err
	}
	return nil
}

func Set(key string, val interface{}, opts ...Option) error {
	options := NewOptions(opts...)
	ttl := time.Duration(0)
	if options.TTL > 0 {
		ttl = options.TTL
	}
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := redis.Ctx(context.Background()).Set(key, b, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func copyInterface(dst interface{}, src interface{}) {
	dstof := reflect.ValueOf(dst)
	if dstof.Kind() == reflect.Ptr {
		el := dstof.Elem()
		srcof := reflect.ValueOf(src)
		if srcof.Kind() == reflect.Ptr {
			el.Set(srcof.Elem())
		} else {
			el.Set(srcof)
		}
	}
}
