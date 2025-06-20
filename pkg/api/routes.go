package api

import (
	"mbti-site/pkg/cache"
	"mbti-site/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine, cacheManager *cache.Manager, i18nManager *i18n.Manager) {
	// 创建处理器
	handler := NewHandler(cacheManager, i18nManager)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 语言路由组
		langGroup := v1.Group("/:lang")
		{
			// 人格类型路由
			typesGroup := langGroup.Group("/types")
			{
				typesGroup.GET("", handler.GetPersonalityTypes)
				typesGroup.GET("/:code/:gender", handler.GetPersonalityType)
			}

			// 人格配对路由
			pairingsGroup := langGroup.Group("/pairings")
			{
				pairingsGroup.GET("", handler.GetPersonalityPairings)
				pairingsGroup.GET("/:male/:female", handler.GetPersonalityPairing)
			}

			// 随机路由
			langGroup.GET("/random/type", handler.GetRandomPersonalityType)
			langGroup.GET("/random/pairing", handler.GetRandomPersonalityPairing)
		}
	}

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// 静态文件服务（如果需要）
	r.Static("/static", "./static")
}
