package gin

import (
	"database/sql"
	"embed"

	"1litw/application"
	"1litw/infrastructure/external"
	"1litw/infrastructure/repository"
	"1litw/presentation/gin/handler"

	"github.com/gin-gonic/gin"
	"github.com/simbafs/kama"
)

func SetupRouter(db *sql.DB, jwtSecret string, webDist embed.FS) *gin.Engine {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewShortURLRepository(db)
	analyticsRepo := repository.NewClickRepository(db)

	// Initialize external services
	uaParser := external.NewUAParserService()

	// Initialize use cases
	userUseCase := application.NewUserUseCase(userRepo, jwtSecret)
	urlUseCase := application.NewURLUseCase(urlRepo, userRepo, analyticsRepo, uaParser)
	analyticsUseCase := application.NewAnalyticsUseCase(analyticsRepo, urlRepo)
	telegramUseCase := application.NewTelegramUseCase(db, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUseCase, telegramUseCase)
	urlHandler := handler.NewURLHandler(urlUseCase, analyticsUseCase)
	userHandler := handler.NewUserHandler(userUseCase)

	// Setup router
	router := gin.Default()
	authed := router.Group("/").Use(handler.AuthMiddleware(jwtSecret, userUseCase))

	// API routes
	// routes about authentication
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/logout", authHandler.Logout)
	authed.POST("/api/auth/telegram/link", authHandler.LinkTelegram)

	// route about user itself
	authed.GET("/api/me", userHandler.GetMe)

	// routes about a short URL
	router.POST("/api/url", handler.OptionalAuthMiddleware(jwtSecret, userUseCase), urlHandler.CreateShortURL)
	authed.GET("/api/url", urlHandler.GetMyURLs)
	authed.DELETE("/api/url/:id", urlHandler.DeleteShortURL)
	authed.GET("/api/url/:id/stats", urlHandler.GetStats)

	// routes about managge users
	authed.GET("/api/user", userHandler.List)
	authed.PUT("/api/user/:id/permission", userHandler.UpdatePermissions)
	authed.DELETE("/api/user/:id", userHandler.Delete)

	// routes about admin
	authed.GET("/api/admin/url", urlHandler.GetAllURLs)

	// Redirection routes
	router.GET("/r/:short_path", urlHandler.Redirect)
	router.GET("/r/@:username/:custom_path", func(c *gin.Context) {
		username := c.Param("username")
		customPath := c.Param("custom_path")
		fullPath := "@" + username + "/" + customPath
		c.Params = append(c.Params, gin.Param{Key: "short_path", Value: fullPath})
		urlHandler.Redirect(c)
	})

	k := kama.New(webDist,
		kama.WithDevServer("http://localhost:4321"),
		// kama.WithTree("/tree"),
		kama.WithPath("web/dist"),
	)

	router.Use(k.Gin())

	return router
}
