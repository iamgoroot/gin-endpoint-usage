package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	ginusagestats "github.com/iamgoroot/gin-usage-stats"
	"github.com/redis/go-redis/v9"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	m := &ginusagestats.StatMiddleware{Backend: &ginusagestats.RedisBackend{RedisClient: rdb}}
	m.Setup(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return r
}

func main() {
	r := setupRouter()
	err := r.Run(":9999")
	fmt.Println(err)
}
