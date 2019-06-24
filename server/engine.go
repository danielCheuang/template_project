package server

import (
	"template_project/router"
	"net/http"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type EngineConfig struct {
	middleware       []func(*gin.Context)
	LimitConnections int
	RunMode          string
	RootRouterPrefix string
}

func (config *EngineConfig) initEngineConfig(api *API) *gin.Engine {
	if config == nil {
		panic("engine config should not be nil")
	}
	conf := api.config
	gin.SetMode(conf.Server.RunMode)
	e := gin.New()
	e.Use(gzip.Gzip(gzip.DefaultCompression))

	e.Use(gin.Logger())
	// use recovery middleware
	e.Use(gin.Recovery())

	e.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTION"},
		AllowHeaders: []string{"Origin"},
		ExposeHeaders: []string{
			"Content-Length",
			"Accept-Language",
			"DNT",
			"X-Mx-ReqToken",
			"Keep-Alive",
			"User-Agent",
			"X-Requested-With",
			"If-Modified-Since",
			"Cache-Control",
			"Content-Type",
			"Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	// By default, http.ListenAndServe (which gin.Run wraps) will serve an unbounded number of requests.
	// Limiting the number of simultaneous connections can sometimes greatly speed things up under load
	if config.LimitConnections > 0 {
		e.Use(limit.MaxAllowed(config.LimitConnections))
	}
	e.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "The incorrect api route",
			"data": nil,
		})
		return
	})

	return e
}

// Init engine init
func (config *EngineConfig) Init(api *API) http.Handler {
	e := config.initEngineConfig(api)
	router.InitRouters(e, *api.config)
	return e
}
