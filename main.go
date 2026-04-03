package main

import (
	"log"
	"os"

	"AI-assistent/config"
	"AI-assistent/handler"
	"AI-assistent/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 日志输出到 stdout，方便 K8s 日志收集
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Gin 切换到 release 模式，减少 debug 日志噪音
	gin.SetMode(gin.ReleaseMode)

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化服务
	ttsService := service.NewTTSService(cfg)
	ttsHandler := handler.NewTTSHandler(ttsService)

	// 创建Gin路由器
	router := gin.Default()

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// 静态页面
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/index.html")
	})

	// 健康检查
	router.GET("/health", ttsHandler.HealthCheck)

	// TTS相关路由
	api := router.Group("/api/v1")
	{
		api.POST("/tts/convert", ttsHandler.ConvertText)
		api.POST("/tts/batch", ttsHandler.BatchConvert)
		api.GET("/tts/voices", ttsHandler.GetVoices)
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
