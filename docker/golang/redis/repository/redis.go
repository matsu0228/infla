package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// Redis is a Redis Driver wrapping go-redis.
type Redis struct {
	client RedisClient
	prefix string
	expire time.Duration
}

// RedisClient is an interface for handling regardless of cluster and single mode.
type RedisClient interface {
	Get(key string) *redis.StringCmd
	HGetAll(key string) *redis.StringStringMapCmd
	Ping() *redis.StatusCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// NewRedis returns instance of Redis.
func NewRedis(isCluster bool, host string, pref string) (*Redis, error) {
	log.Printf("[DEBUG] NewRedis() with isCluster:%v, host:%v, prefix:%v", isCluster, host, pref)
	var client RedisClient
	if isCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{host},
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr: host,
		})
	}
	// PINGで疎通確認
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return &Redis{
		client: client,
		prefix: pref,
		expire: 24 * time.Hour, //TODO: 外から指定できるように
	}, nil
}

func (r *Redis) key(pID string) string {
	return fmt.Sprintf("%s:%s", r.prefix, pID)
}

// Get : Redisから"r.prefix+pID"をキーとして値を取得する
func (r *Redis) Get(pID string) (string, error) {
	key := r.key(pID)
	log.Printf("[DEBUG] Get() with %v", key)
	// result, err := r.client.HGetAll(key).Result()
	// if err != nil {
	// 	return "", err
	// }
	// return result["Comment"], nil

	return r.client.Get(key).Result()
}

// Save :Redisから"r.prefix+pID"をキーとして値を保存
func (r *Redis) Save(pID, value string) error {

	key := r.key(pID)
	log.Printf("[DEBUG] Save() with key:%v, value:%v", key, value)
	// m := myredisi.StructToMap(c)
	// 値と期限のセットをアトミックに処理
	// pipe := r.Pipeline()
	// pipe.HMSet(k, v)
	// pipe.Expire(k, redisExpire)
	// _, err := pipe.Exec()
	// if err != nil {
	// 	return err
	// }

	return r.client.Set(key, value, r.expire).Err()
}
