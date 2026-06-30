package app

import (
	"mouban/internal/controller"
	"mouban/internal/util"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	serviceOpsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "mouban_service_ops_duration",
		Help:    "Histogram of the duration of HTTP service requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "ua", "referer"})
)

func RunHTTPServer() error {
	router := NewRouter()
	return router.Run(":" + viper.GetString("server.port"))
}

func NewRouter() *gin.Engine {
	mode := strings.ToLower(strings.TrimSpace(viper.GetString("server.mode")))
	switch mode {
	case "", gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.DebugMode, gin.TestMode:
		gin.SetMode(mode)
	default:
		logrus.Warnln("invalid server.mode:", mode, "fallback to release")
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Static("/assets", "build/assets")
	router.StaticFile("/favicon.ico", "build/favicon.ico")
	router.StaticFile("/robots.txt", "build/robots.txt")

	router.Use(recoverMiddleware)
	router.Use(corsMiddleware)
	router.Use(accessLogMiddleware)
	router.Use(metricsMiddleware)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/swagger/")
	})
	router.Static("/swagger", "build/swagger")

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/refresh_item", controller.RefreshItem)
		adminGroup.GET("/refresh_user", controller.RefreshUser)
	}

	queryGroup := router.Group("/guest")
	{
		queryGroup.GET("/resolve_user", controller.ResolveUser)
		queryGroup.GET("/check_user", controller.CheckUser)
		queryGroup.GET("/item_detail", controller.GuestItemDetail)
		queryGroup.GET("/user_book", controller.GuestUserBook)
		queryGroup.GET("/user_game", controller.GuestUserGame)
		queryGroup.GET("/user_movie", controller.GuestUserMovie)
		queryGroup.GET("/user_song", controller.GuestUserSong)
	}

	router.GET("/explore/queue_overview", controller.QueueOverview)

	registerSPAFallback(router)

	return router
}

func registerSPAFallback(router *gin.Engine) {
	router.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.JSON(http.StatusNotFound, gin.H{"success": false, "msg": "not found"})
			return
		}

		requestPath := ctx.Request.URL.Path
		if strings.HasPrefix(requestPath, "/guest/") ||
			strings.HasPrefix(requestPath, "/admin/") ||
			requestPath == "/metrics" ||
			requestPath == "/explore/queue_overview" {
			ctx.JSON(http.StatusNotFound, gin.H{"success": false, "msg": "not found"})
			return
		}

		if strings.Contains(path.Base(requestPath), ".") {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.File("build/index.html")
	})
}

func recoverMiddleware(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"msg":     "服务内部错误，请联系开发者处理",
			})
			logrus.Errorln("handle panic", r, "=>", util.GetCurrentGoroutineStack())
		}
	}()
	ctx.Next()
}

func corsMiddleware(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")

	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

	c.Next()
}

func accessLogMiddleware(c *gin.Context) {
	startTime := time.Now()
	c.Next()

	latencyTime := time.Since(startTime)
	reqURI := c.Request.RequestURI
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()

	logrus.Infoln("uri", reqURI, "status_code", statusCode, "cost", latencyTime, "client_ip", clientIP)
}

func metricsMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	if c.Writer.Status() == http.StatusNotFound {
		return
	}

	if c.Request.RequestURI == "/metrics" {
		return
	}

	referer := c.Request.Referer()
	if referer == "" {
		referer = "-"
	}

	duration := time.Since(start).Seconds()
	serviceOpsHistogram.WithLabelValues(c.Request.Method, c.Request.URL.Path, c.Request.UserAgent(), referer).Observe(duration)
}
