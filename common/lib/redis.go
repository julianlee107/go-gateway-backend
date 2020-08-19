package lib

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func RedisConnFactory(name string) (redis.Conn, error) {
	if ConfRedisMap != nil && ConfRedisMap.Map != nil {
		if conf, ok := ConfRedisMap.Map[name]; ok {
			randHost := conf.ProxyList[rand.Intn(len(conf.ProxyList))]
			if conf.WriteTimeout == 0 {
				conf.WriteTimeout = 50
			}
			if conf.ReadTimeout == 0 {
				conf.ReadTimeout = 50
			}
			if conf.ConnTimeout == 0 {
				conf.ConnTimeout = 50
			}

			conn, err := redis.Dial(
				"tcp",
				randHost,
				redis.DialConnectTimeout(time.Duration(conf.ConnTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(conf.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(conf.WriteTimeout)*time.Millisecond),
			)
			if err != nil {
				return nil, err
			}

			if conf.Password != "" {
				if _, err := conn.Do("AUTH", conf.Password); err != nil {
					conn.Close()
					return nil, err
				}

				if conf.Db != 0 {
					if _, err := conn.Do("SELECT", conf.Db); err != nil {
						conn.Close()
						return nil, err
					}
				}
				return conn, nil
			}
		}
	}
	return nil, errors.New("create redis conn fail")
}

func RedisConfDo(trace *TraceContext, name string, commandName string, args ...interface{}) (interface{}, error) {
	conn, err := RedisConnFactory(name)
	if err != nil {
		Log.TagError(trace, "redis_failure", map[string]interface{}{
			"method": commandName,
			"err":    errors.New("RedisConnFactory error: " + name),
			"bind":   args,
		})
		return nil, err
	}

	defer conn.Close()

	startExecTime := time.Now()
	reply, err := conn.Do(commandName, args)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		replyStr, _ := redis.String(reply, nil)
		Log.TagError(trace, "redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"reply":     replyStr,
			"proc_time": fmt.Sprintf("%fs", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, nil
}
