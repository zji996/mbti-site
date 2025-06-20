package main

import (
	"flag"
	"log"

	"mbti-site/pkg/api"
	"mbti-site/pkg/cache"
	"mbti-site/pkg/content"
	"mbti-site/pkg/i18n"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	var (
		port        = flag.String("port", "8080", "Server port")
		contentPath = flag.String("content", "./content", "Content directory path")
		debug       = flag.Bool("debug", false, "Enable debug mode")
	)
	flag.Parse()

	// 设置Gin模式
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化内容解析器
	contentParser := content.NewParser(*contentPath)

	// 初始化缓存
	cacheManager := cache.NewManager(contentParser)

	// 初始化国际化
	i18nManager := i18n.NewManager()

	// 创建Gin引擎
	r := gin.Default()

	// 设置中间件
	setupMiddleware(r)

	// 设置API路由
	api.SetupRoutes(r, cacheManager, i18nManager)

	// 启动服务器
	log.Printf("Starting server on port %s", *port)
	if err := r.Run(":" + *port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupMiddleware(r *gin.Engine) {
	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 日志中间件
	r.Use(gin.Logger())

	// 恢复中间件
	r.Use(gin.Recovery())
}
