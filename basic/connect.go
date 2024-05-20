package main

import "github.com/redis/go-redis/v9"

func RedisConnect() *redis.Client {
	url := "redis://localhost:6379"
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opts)
}
