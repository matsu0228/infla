package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/matsu0228/infla/redis/repository"
)

// errorExit :エラー終了時の共通処理
func errorExit(err error) {
	log.Fatal("[ERROR] ", err)
}

func envLoad() error {
	if err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))); err != nil {
		return err
	}
	return nil
}

func isCluster(mode string) bool {
	switch strings.ToLower(mode) {
	case "cluster":
		return true
		// case "single":
		// default:
	}
	return false
}

func main() {

	if err := envLoad(); err != nil {
		errorExit(fmt.Errorf("can not load .env.%s :err=%v", os.Getenv("APP_ENV"), err))
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClusterMode := os.Getenv("REDIS_CLUSTER_MODE")
	redisPrefix := os.Getenv("REDIS_PREFIX")

	redis, err := repository.NewRedis(isCluster(redisClusterMode), redisAddr, redisPrefix)
	if err != nil {
		errorExit(err)
	}

	// pp.Print(redis)
	err = redis.Save("key02", "value002")
	if err != nil {
		errorExit(err)
	}
	value, err := redis.Get("key02")
	if err != nil {
		errorExit(err)
	}
	fmt.Println("redis VALUE=", value)

}
