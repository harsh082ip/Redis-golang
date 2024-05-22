package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/go-redis/redis/v9"
	connect "github.com/harsh082ip/Redis-golang/Connect"
	"github.com/redis/go-redis/v9"
)

const (
	WEBPORT = ":8002"
)

type User struct {
	Name     string `json:"name"`
	Job      string `json:"job"`
	Location string `json:"location"`
	IP       string `json:"ip"`
}

func main() {
	router := gin.Default()
	rdb := connect.RedisConnect()

	router.GET("/getdetails", func(ctx *gin.Context) {
		var user User

		r, err := rdb.Get(ctx.Request.Context(), "data:user1").Result()
		if err == redis.Nil {
			log.Println("Cache miss, fetching data...")
		} else if err != nil {
			log.Println("Error getting data from Redis:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data from cache"})
			return
		} else {
			err = json.Unmarshal([]byte(r), &user)
			if err != nil {
				log.Println("Error decoding the struct:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode data"})
				return
			}
			ctx.JSON(http.StatusOK, user)
			return
		}

		// Simulate fetching from a data source
		time.Sleep(time.Second * 8)
		user = User{
			Name:     "Harsh Vardhan Singh",
			Job:      "Software Developer",
			Location: "UK, London",
			IP:       "172.28.8.1",
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Println("Error marshalling the data:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
			return
		}

		_, err = rdb.Set(ctx.Request.Context(), "data:user1", jsonData, time.Second*120).Result()
		if err != nil {
			log.Println("Error setting data in Redis:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set data in cache"})
			return
		}
		log.Println("Data cached successfully")

		ctx.JSON(http.StatusOK, user)
	})

	router.POST("/updatedetails", func(ctx *gin.Context) {
		var user User

		if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
			log.Println("Error decoding request body:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		existingData, err := rdb.Get(ctx.Request.Context(), "data:user1").Result()
		if err != nil && err == redis.Nil {
			log.Println("Error getting data from Redis:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data from cache"})
			return
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			log.Println("Error marshalling data:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
			return
		}

		if existingData != string(jsonData) {
			err := rdb.Del(ctx.Request.Context(), "data:user1").Err()
			if err != nil {
				log.Println("Error deleting data from Redis:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cache"})
				return
			}

			// Simulate updating in the database
			log.Println("Data updated in the database")

			ctx.JSON(http.StatusOK, gin.H{"status": "Updation successful"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Data is not changed"})
	})

	if err := http.ListenAndServe(WEBPORT, router); err != nil {
		log.Println("Error starting the server:", err)
	}
}
