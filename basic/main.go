package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379", // Redis address
	// 	Password: "",               // No password set
	// 	DB:       0,                // Use default DB
	// })
	// opt, err := redis.ParseURL("redis://localhost:6379")
	// if err != nil {
	// 	panic(err)
	// }

	// rdb := redis.NewClient(opt)

	rdb := RedisConnect()
	// fmt.Println(rdb)

	// Ping the Redis server
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)

	val, _ := rdb.Set(ctx, "user:1", "harsh", 0).Result()
	fmt.Println(val)
	// fmt.Println(rdb)

	res, err := rdb.Del(ctx, "user:1").Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)

	success, err := rdb.SetNX(ctx, "user:1", "ankiit", 0).Result()
	if err != nil {
		fmt.Println("error")
		fmt.Println(success, err.Error())
	}
	fmt.Println(success)
	val, _ = rdb.Get(ctx, "name").Result()
	fmt.Println(val)
}
