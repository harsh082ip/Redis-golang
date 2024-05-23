package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	connect "github.com/harsh082ip/Redis-golang/Connect"
	"github.com/redis/go-redis/v9"
)

/*

Problem 5: Implementing a Session Store

Create a session management system using Redis to store user session data.

Requirements:

    Use Gin for the web framework.
    Implement middleware to handle session creation and validation.
    Store session data in Redis with a timeout (e.g., 30 minutes).
    Provide endpoints for login, logout, and fetching session data.
    Ensure that sessions are invalidated upon logout.

*/

const (
	WEBPORT     = ":8002"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionTTL  = 30 * time.Second
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionInfo struct {
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
}

func main() {

	router := gin.Default()

	var users []User

	rdb := connect.RedisConnect()

	router.POST("/signup", func(ctx *gin.Context) {

		var user User

		if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Error in Decoding the request body",
			})
			return
		}

		users = append(users, user)

		ctx.JSON(http.StatusOK, gin.H{
			"msg": "SignUp Successful",
		})

	})

	router.POST("/login", func(ctx *gin.Context) {

		var user LoginUser

		if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Error in Decoding the request body",
			})
			return
		}

		isPresent := false

		for _, u := range users {

			if u.Email == user.Email && u.Password == user.Password {
				isPresent = true
				break
			}
		}

		if !isPresent {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Entered email and password is not present",
			})
			return
		}

		session_id, err := CreateSeessionID(30)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Error in Creating the SessionID",
			})
		}

		session_info := &SessionInfo{
			SessionID: session_id,
			Email:     user.Email,
		}
		key := "session_info:" + session_id

		jsonData, _ := json.Marshal(session_info)

		res, err := rdb.Set(context.TODO(), key, jsonData, sessionTTL).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error in Registering the session info",
				"msg":   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":      res,
			"session_id":  session_info.SessionID,
			"expiring in": "120 seconds",
		})
	})

	router.POST("/test", func(ctx *gin.Context) {

		var session_info SessionInfo

		if err := json.NewDecoder(ctx.Request.Body).Decode(&session_info); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "Error in Decoding the request body",
			})
			return
		}

		key := "session_info:" + session_info.SessionID

		_, err := rdb.Get(context.TODO(), key).Result()
		if err != nil {
			if err == redis.Nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized Access",
					"msg":   "No such session Id is present",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Authorized to access, test successful",
		})
		// return
	})

	router.POST("/logout", func(ctx *gin.Context) {
		var sessionInfo SessionInfo
		if err := ctx.ShouldBindJSON(&sessionInfo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Error in decoding the request body"})
			return
		}

		key := "session_info:" + sessionInfo.SessionID
		if _, err := rdb.Del(context.TODO(), key).Result(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error in deleting the session", "msg": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "Logout successful"})
	})

	if err := http.ListenAndServe(WEBPORT, router); err != nil {
		log.Fatal("Error starting the server, ", err.Error())
	}
}

func CreateSeessionID(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		result[i] = letterBytes[num.Int64()]
	}
	return string(result), nil
}
