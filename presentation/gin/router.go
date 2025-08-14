package gin

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"1litw/application"
	"1litw/infrastructure/external"
	"1litw/infrastructure/repository"
	"1litw/presentation/gin/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, jwtSecret string, webDist embed.FS) *gin.Engine {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewShortURLRepository(db)
	analyticsRepo := repository.NewClickRepository(db)

	// Initialize external services
	uaParser := external.NewUAParserService()
	geoIP := external.NewGeoIPService()

	// Initialize use cases
	userUseCase := application.NewUserUseCase(userRepo, jwtSecret)
	urlUseCase := application.NewURLUseCase(urlRepo, userRepo, analyticsRepo, uaParser, geoIP)
	analyticsUseCase := application.NewAnalyticsUseCase(analyticsRepo, urlRepo)
	telegramUseCase := application.NewTelegramUseCase(db, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUseCase, telegramUseCase)
	urlHandler := handler.NewURLHandler(urlUseCase, analyticsUseCase)

	// Setup router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			// Telegram linking endpoint is authenticated
		}

		// Authenticated routes
		authRequired := api.Group("/")
		authRequired.Use(handler.AuthMiddleware(jwtSecret, userUseCase))
		{
			authRequired.POST("/auth/telegram/link", authHandler.LinkTelegram)

			// urls.POST("", urlHandler.CreateShortURL) // Will be moved
			authRequired.GET("", urlHandler.GetMyURLs)
			authRequired.DELETE("/:id", urlHandler.DeleteShortURL)
			authRequired.GET("/:id/stats", urlHandler.GetStats)
			// Admin routes would be nested here with another middleware
		}

		// URL creation can be done by anonymous users
		api.POST("/urls", handler.OptionalAuthMiddleware(jwtSecret, userUseCase), urlHandler.CreateShortURL)
	}

	// Redirection routes
	router.GET("/:short_path", urlHandler.Redirect)
	router.GET("/@:username/:custom_path", func(c *gin.Context) {
		username := c.Param("username")
		customPath := c.Param("custom_path")
		fullPath := "@" + username + "/" + customPath
		c.Params = append(c.Params, gin.Param{Key: "short_path", Value: fullPath})
		urlHandler.Redirect(c)
	})

	// Serve frontend
	distFS, err := fs.Sub(webDist, "web/dist")
	if err != nil {
		panic(err)
	}

	// Serve static assets from the /assets directory
	router.StaticFS("/assets", http.FS(distFS))

	// For all other routes, serve the index.html, letting the client-side router take over.
	// This is for SPAs. Since Astro can be an MPA, we need to be more specific.
	router.NoRoute(func(c *gin.Context) {
		// If it's not an API route and not a file, serve the corresponding html file or the index.html
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			// Try to serve the file directly from the embed FS
			// e.g. /dashboard -> /dashboard/index.html
			filePath := c.Request.URL.Path
			if strings.HasSuffix(filePath, "/") {
				filePath = filePath + "index.html"
			} else if !strings.Contains(filePath, ".") {
				// This is a heuristic. If there's no dot, it's likely a page route.
				// Astro builds pages as /pagename/index.html
				filePath = filePath + "/index.html"
			}

			// Serve the file if it exists
			if _, err := distFS.Open(strings.TrimPrefix(filePath, "/")); err == nil {
				c.FileFromFS(filePath, http.FS(distFS))
				return
			}

			// Fallback to root index.html for true SPAs or as a default
			c.FileFromFS("/", http.FS(distFS))
		}
		// If it's an API route that doesn't match, Gin will handle the 404
	})

	return router
}
