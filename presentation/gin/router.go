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

	// API routes
	api := router.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/logout", authHandler.Logout)
		// Telegram linking endpoint is authenticated

		// Authenticated routes
		{
			authRequired := api.Group("/").Use(handler.AuthMiddleware(jwtSecret, userUseCase))

			authRequired.GET("/me", userHandler.GetMe)

			authRequired.POST("/auth/telegram/link", authHandler.LinkTelegram)

			authRequired.GET("/url", urlHandler.GetMyURLs)
			authRequired.DELETE("/url/:id", urlHandler.DeleteShortURL)
			authRequired.GET("/url/:id/stats", urlHandler.GetStats)

			authRequired.GET("/user", userHandler.List)
			authRequired.PUT("/user/:id/permission", userHandler.UpdatePermissions)
			authRequired.DELETE("/user/:id", userHandler.Delete)

			authRequired.GET("/admin/urls", urlHandler.GetAllURLs)
		}

		// URL creation can be done by anonymous users
		api.POST("/url", handler.OptionalAuthMiddleware(jwtSecret, userUseCase), urlHandler.CreateShortURL)
	}

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
