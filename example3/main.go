package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	connect "github.com/harsh082ip/Redis-golang/Connect"
	"github.com/redis/go-redis/v9"
)

/*

Problem 4: Distributed Locking

Implement a distributed locking mechanism using Redis to ensure that a critical section of code is executed by only one instance in a distributed system.

Requirements:

    Use Gin for the web framework.
    Use Redis to implement a distributed lock (e.g., using the SET command with NX and EX options).
    Create an endpoint that simulates a critical operation (e.g., updating a shared resource).
    Ensure that the critical operation is only performed by one instance at a time across multiple instances.

*/

const (
	WEBPORT = ":8002"
)

type Lock struct {
	Owner      string `json:"owner"`
	Lock_value string `json:"lock_val"`
}

func main() {

	router := gin.Default()

	router.POST("/setlock", SetLock)
	router.GET("/getlockinfo", GetLockInfo)

	if err := http.ListenAndServe(WEBPORT, router); err != nil {
		log.Fatal("Error Starting the port: ", err.Error())
	}
}

func GetLockInfo(c *gin.Context) {

	rdb := connect.RedisConnect()

	var lock Lock

	res, err := rdb.Get(context.Background(), "lockey").Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "No lock available currently",
			})
			return
		}
		log.Fatal("Error in Getting val from Redis, ", err.Error())
	}

	if err := json.Unmarshal([]byte(res), &lock); err != nil {
		log.Fatal("Error in Unmarshalling the data")
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Lock Present",
		"info":   lock,
	})
}

func SetLock(c *gin.Context) {

	rdb := connect.RedisConnect()

	var lock Lock

	if err := json.NewDecoder(c.Request.Body).Decode(&lock); err != nil {
		log.Fatal("Error Decoding the request body", err.Error())
	}

	jsonData, err := json.Marshal(lock)
	if err != nil {
		log.Fatal("Error Marshalling the data...")
	}

	accquired, _ := rdb.SetNX(context.Background(), "lockey", jsonData, time.Second*10).Result()

	if !accquired {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Try again after some time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "lock sucessfully set",
	})
	// return

}
