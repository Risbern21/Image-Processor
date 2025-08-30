package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

func Client() *redis.Client {
	return cache
}

func Connect() {
	addr := os.Getenv("REDIS_ADDR")
	passwd := os.Getenv("REDIS_PASSWD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	ctx := context.Background()

	cmd := rdb.Ping(ctx)
	if cmd.Err() != nil {
		log.Printf("unable to connect to cache: %v", cmd.Err())
	}

	log.Println("successfully connected to the cache")

	cache = rdb
}
