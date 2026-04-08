package main

import (
	"log"
	"os"

	"github.com/liyufan816-web/mimo-voice-assistant/config"
	"github.com/liyufan816-web/mimo-voice-assistant/handler"
	"github.com/liyufan816-web/mimo-voice-assistant/service"

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
	appCfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载应用配置失败：%v", err)
	}

	// 初始化数据库连接
	dbCfg := config.LoadDatabaseConfig()
	if err := config.InitDB(dbCfg); err != nil {
		log.Fatalf("初始化数据库失败：%v", err)
	}
	defer func() {
		if err := config.CloseDB(); err != nil {
			log.Printf("关闭数据库连接失败：%v", err)
		}
	}()

	// 初始化服务
	ttsService := service.NewTTSService(appCfg)
	ttsHandler := handler.NewTTSHandler(ttsService)
	authHandler := handler.NewAuthHandler()
	historyHandler := handler.NewHistoryHandler()

	// 创建 Gin 路由器
	router := gin.Default()

	// 配置 CORS
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

	// 公开路由（无需认证）
	public := router.Group("/api/v1")
	{
		// 认证相关
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/auth/check-exists", authHandler.CheckUserExists)
	}

	// 需要认证的路由
	api := router.Group("/api/v1")
	api.Use(authHandler.Middleware())
	{
		// TTS 相关
		api.POST("/tts/convert", ttsHandler.ConvertText)
		api.POST("/tts/batch", ttsHandler.BatchConvert)
		api.GET("/tts/voices", ttsHandler.GetVoices)

		// 用户信息
		api.GET("/user/info", authHandler.GetUserInfo)
		api.GET("/user/api-key", authHandler.GetAPIKey)

		// 历史记录
		api.GET("/history", historyHandler.GetUserHistory)
		api.GET("/history/:id", historyHandler.GetRecordDetail)
		api.DELETE("/history/:id", historyHandler.DeleteRecord)
		api.DELETE("/history", historyHandler.ClearUserHistory)
		api.GET("/history/export", historyHandler.ExportUserHistory)
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", appCfg.ServerPort)
	if err := router.Run(":" + appCfg.ServerPort); err != nil {
		log.Fatalf("启动服务器失败：%v", err)
	}
}
