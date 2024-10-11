package ginusagestats

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type Stat struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Count    int64  `json:"count"`
}

type StatMiddleware struct {
	Backend interface {
		Collect(ctx context.Context, method, endpoint string, incr int64) error
		GetStats(ctx context.Context) ([]Stat, error)
	}
}

func (m *StatMiddleware) Setup(router *gin.Engine) {
	router.GET("/endpoint-usage-stats/:type", m.StatsHandler)
	router.GET("/endpoint-usage-stats", m.StatsHandler)
	router.Use(m.Stat)
	routes := router.Routes()
	for _, route := range routes {
		if strings.HasPrefix(route.Path, "/endpoint-usage-stats") { // skip own route
			continue
		}
		_ = m.Backend.Collect(context.Background(), route.Method, route.Path, 0)
	}
}

func (m *StatMiddleware) Stat(ctx *gin.Context) {
	method := ctx.Request.Method
	endpoint := ctx.FullPath()
	if endpoint == "" {
		ctx.Next()
		return
	}
	err := m.Backend.Collect(ctx.Request.Context(), method, endpoint, 1)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Next()
}

func (m *StatMiddleware) StatsHandler(ctx *gin.Context) {
	stats, err := m.Backend.GetStats(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slices.SortFunc(stats, sortStatsFunc)
	contentType := ctx.Params.ByName("type")
	switch contentType {
	case "json":
		ctx.JSON(http.StatusOK, stats)
	case "csv":
		ctx.Header("Content-Type", "text/csv")
		ctx.Header("Content-Disposition", "attachment; filename=endpoint-usage-stats.csv")
		for _, stat := range stats {
			fmt.Fprintf(ctx.Writer, "%s,%s,%d\n", stat.Method, stat.Endpoint, stat.Count)
		}
	case "xml":
		ctx.XML(http.StatusOK, stats)
	default:
		err := tableTemplate.Execute(ctx.Writer, stats)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

func sortStatsFunc(a, b Stat) int {
	if a.Count == b.Count {
		return strings.Compare(a.Endpoint, b.Endpoint)
	}
	if a.Count == b.Count {
		return 0
	}
	if a.Count > b.Count {
		return -1
	}
	return 1
}
