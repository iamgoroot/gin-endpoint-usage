package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	ginendpointusage "github.com/iamgoroot/gin-endpoint-usage"
	"github.com/redis/go-redis/v9"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	m := &ginendpointusage.StatMiddleware{Backend: &ginendpointusage.RedisBackend{RedisClient: rdb}}
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
