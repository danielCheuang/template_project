package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Host        string
	Port        string
	Password    string
	IdleTimeout int
	MaxIdle     int
	MaxActive   int
}

type Service struct {
	config Config
	pool   *redis.Pool
}

// Initialize redis service init
func (service *Service) Initialize(config interface{}) {
	cfg, ok := config.(Config)
	if !ok {
		panic("redis service config error!")
	}

	service.config = cfg
	service.pool = &redis.Pool{
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
			var (
				c   redis.Conn
				err error
			)
			if c, err = redis.Dial("tcp", address); err != nil {
				return nil, err
			}
			if len(cfg.Password) > 0 {
				if _, err = c.Do("AUTH", cfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
}

// Start xx
func (service *Service) Start() {

}

// Stop xx
func (service *Service) Stop() {

}

// -----------------string operation------------------
// when set exist key, old key ttl must reset it
func (service *Service) Set(key string, value []byte, ttl int64) error {
	conn := service.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("set", key, value); err != nil {
		return err
	}

	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {
			return err
		}
	}
	return nil
}

func (service *Service) Get(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("get", key)
	if nil != err {
		return []byte{}, err
	} else if nil == reply {
		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *Service) Mset(keyvalues [][2][]byte, ttl int64) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(keyvalues) <= 0 {
		return nil
	}

	var list []interface{}
	for _, v := range keyvalues {
		list = append(list, v[0])
		list = append(list, v[1])
	}

	if _, err := conn.Do("mset", list...); err != nil {
		return err
	}

	if ttl > 0 {
		for _, kv := range keyvalues {
			if _, err := conn.Do("expire", kv[0], ttl); err != nil {
				return err
			}
		}
	}
	return nil
}

func (service *Service) Mget(keys []string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	if keys == nil || len(keys) <= 0 {
		return [][]byte{}, nil
	}
	var list []interface{}
	for _, v := range keys {
		list = append(list, v)
	}
	reply, err := conn.Do("mget", list...)
	res := [][]byte{}
	if nil != err {

		return [][]byte{}, err
	} else if nil == reply {

		return [][]byte{}, err
	} else {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil != r {
				res = append(res, r.([]byte))
			} else {
				res = append(res, nil)
			}
		}
		return res, err
	}
}

func (service *Service) Exists(key string) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("exists", key)
	if err != nil {

		return false, err
	} else {
		exists, _ := reply.(int64)
		if exists == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (service *Service) Del(key string) error {
	conn := service.pool.Get()
	defer conn.Close()

	_, err := conn.Do("del", key)
	if nil != err {
		return err
	}
	return err
}

func (service *Service) Dels(keys []string) error {
	conn := service.pool.Get()
	defer conn.Close()

	if keys == nil || len(keys) <= 0 {
		return nil
	}

	var list []interface{}
	for _, v := range keys {
		list = append(list, v)
	}

	_, err := conn.Do("del", list...)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) Keys(keyFormat string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("keys", keyFormat)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil != r {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

// -----------------set operation---------------------
func (service *Service) SAdd(key string, ttl int64, members ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range members {
		vs = append(vs, v)
	}
	_, err := conn.Do("sadd", vs...)
	if nil != err {
		return err
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {

			return err
		}
	}
	return err
}

func (service *Service) SRem(key string, members ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range members {
		vs = append(vs, v)
	}
	reply, err := conn.Do("srem", vs...)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) SetNX(key string, value int64, ttl int64) error {
	conn := service.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("setnx", key, value); err != nil {

		return err
	}

	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {

			return err
		}
	}
	return nil
}

func (service *Service) SCard(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("scard", key)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) SIsMember(key string, member []byte) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("sismember", key, member)
	if err != nil {

		return false, err
	} else {
		return reply.(int64) > 0, nil
	}
}

func (service *Service) SMembers(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("smembers", key)

	res := [][]byte{}
	if nil != err {

	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

// -----------------zset operation-------------------
func (service *Service) ZAdd(key string, ttl int64, args ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(args)%2 != 0 {
		return errors.New("the length of `args` must be even")
	}
	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	_, err := conn.Do("zadd", vs...)
	if nil != err {
		return err
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {

			return err
		}
	}
	return err
}

func (service *Service) ZRem(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("zrem", vs...)
	if nil != err {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) ZCard(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("zcard", key)
	if nil != err {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) ZRank(key string, member []byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, member)

	reply, err := conn.Do("zrank", vs...)
	if nil != err {

		return -1, err
	} else if reply == nil {
		return -1, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) ZRevRank(key string, member []byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, member)

	reply, err := conn.Do("zrevrank", vs...)
	if nil != err {

		return -1, err
	} else if reply == nil {
		return -1, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) ZRange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrange", vs...)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *Service) ZRangeByScore(key string, min, max interface{}, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, min, max)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrangebyscore", vs...)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *Service) ZRevRange(key string, start, stop int64, withScores bool) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)
	if withScores {
		vs = append(vs, []byte("WITHSCORES"))
	}
	reply, err := conn.Do("zrevrange", vs...)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *Service) ZRemRangeByScore(key string, start, stop int64) (int64, error) {

	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key, start, stop)

	reply, err := conn.Do("ZREMRANGEBYSCORE", vs...)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

// -----------------hash operation-------------------
func (service *Service) HSet(key string, ttl int64, field string, value []byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, field)
	vs = append(vs, value)

	_, err := conn.Do("hset", vs...)
	if nil != err {
		return err
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {

			return err
		}
	}
	return err
}

func (service *Service) HGet(key string, field []byte) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	var vs []interface{}
	vs = append(vs, key)
	vs = append(vs, field)

	reply, err := conn.Do("hget", vs...)

	if nil != err {

		return []byte{}, err
	} else if nil == reply {

		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *Service) HMSet(key string, ttl int64, args ...[]byte) error {
	conn := service.pool.Get()
	defer conn.Close()

	if len(args)%2 != 0 {
		return errors.New("the length of `args` must be even")
	}
	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	_, err := conn.Do("hmset", vs...)
	if nil != err {
		return err
	}
	if ttl > 0 {
		if _, err := conn.Do("expire", key, ttl); err != nil {

			return err
		}
	}
	return err
}

func (service *Service) HMGet(key string, fields ...[]byte) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range fields {
		vs = append(vs, v)
	}
	reply, err := conn.Do("hmget", vs...)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil == err && nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			if nil == r {
				res = append(res, []byte{})
			} else {
				res = append(res, r.([]byte))
			}
		}
	}
	return res, err
}

func (service *Service) HDel(key string, fields ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range fields {
		vs = append(vs, v)
	}
	reply, err := conn.Do("hdel", vs...)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) HExists(key string, field []byte) (bool, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hexists", key, field)
	if nil != err {
		return false, err
	} else if nil == err && nil != reply {
		exists := reply.(int64)
		return exists > 0, nil
	}

	return false, err
}

func (service *Service) HKeys(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hkeys", key)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *Service) HVals(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hvals", key)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *Service) HGetAll(key string) ([][]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", key)

	res := [][]byte{}
	if nil != err {
		return res, err
	} else if nil != reply {
		rs := reply.([]interface{})
		for _, r := range rs {
			res = append(res, r.([]byte))
		}
	}
	return res, err
}

func (service *Service) HLen(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("hlen", key)

	if nil != err {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

// -----------------list operation--------------------
func (service *Service) LRpush(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("lpush", vs...)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) LLpush(key string, args ...[]byte) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	for _, v := range args {
		vs = append(vs, v)
	}
	reply, err := conn.Do("rpush", vs...)

	if err != nil {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}

func (service *Service) LRpop(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("rpop", key)

	if nil != err {

		return []byte{}, err
	} else if nil == reply {

		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *Service) LLpop(key string) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("lpop", key)

	if nil != err {

		return []byte{}, err
	} else if nil == reply {

		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *Service) LIndex(key string, index int64) ([]byte, error) {
	conn := service.pool.Get()
	defer conn.Close()

	vs := []interface{}{}
	vs = append(vs, key)
	vs = append(vs, index)

	reply, err := conn.Do("lindex", vs)

	if nil != err {

		return []byte{}, err
	} else if nil == reply {

		return []byte{}, err
	} else {
		return reply.([]byte), err
	}
}

func (service *Service) LLlen(key string) (int64, error) {
	conn := service.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("llen", key)
	if nil != err {

		return 0, err
	} else {
		res := reply.(int64)
		return res, err
	}
}
