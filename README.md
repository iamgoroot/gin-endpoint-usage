# gin-usage-stats

A simple middleware for Gin that allows you to collect usage statistics for your endpoints to find unused endpoints.
It's intended to use this on your test server while running your autotests.

For production monitoring use prometheus libs for example https://github.com/penglongli/gin-metrics 

## Installation

You can install the package via go get:

```golang
go get github.com/iamgoroot/gin-endpoint-usage
```

## Usage

You can setup it directly in your Gin application as follows:
```golang
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/iamgoroot/gin-endpoint-usage"
)

func main() {
   	router := gin.Default()
	rdb := redis.NewClient(&redis.Options{ /* redis options */	})
	stats := &StatMiddleware{Backend: &RedisBackend{RedisClient: rdb}}
	stats.Setup(router)

    router.Run(":8080")
}
```

get the stats via endpoints

`http://localhost:8080/endpoint-usage-stats` for html table
`http://localhost:8080/endpoint-usage-stats/json` for json
`http://localhost:8080/endpoint-usage-stats/csv` for csv
`http://localhost:8080/endpoint-usage-stats/xml` for xml