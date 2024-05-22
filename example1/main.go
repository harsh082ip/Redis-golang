package main

/*
****************************************************

 :) An EXAMPLE OF SERVER SIDE CACHING

****************************************************
*/

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

type Response struct {
	Name     string `json:"Name"`
	Job      string `json:"Job"`
	Location string `json:"Location"`
	IP       string `json:"IP"`
}

func main() {

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		rdb := connect.RedisConnect()

		var res Response

		r, err := rdb.Get(context.TODO(), "data:1").Result()
		if err != nil {
			if err == redis.Nil {
				log.Println("Error: No Data Found...")
			}
			log.Println(err)
		} else {
			err := json.Unmarshal([]byte(r), &res)
			if err != nil {
				log.Println(err.Error())
			}

			ctx.JSON(http.StatusOK, res)
			return
		}

		time.Sleep(time.Second * 8)

		res = Response{
			Name:     "Harsh Vardhan Singh",
			Job:      "Software Developer",
			Location: "UK",
			IP:       "192.168.165.207",
		}

		jsonData, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = rdb.Set(context.TODO(), "data:1", jsonData, 120*time.Second).Result()
		if err != nil {
			log.Println(err)
		}

		ctx.JSON(http.StatusOK, res)
		// return
	})

	router.Run(":8002")
	// http.ListenAndServe(":8002", router)
}
