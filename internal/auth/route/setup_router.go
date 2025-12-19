package route

import (
    "github.com/TimX-21/auth-service-go/internal/auth/handler"
    "github.com/TimX-21/auth-service-go/internal/middleware"
    "github.com/gin-gonic/gin"
)

type RouteConfig struct {
    AuthHandler *handler.AuthHandler
}

func NewRouteConfig(
    authH *handler.AuthHandler,
) *RouteConfig {
    return &RouteConfig{
        AuthHandler: authH,
    }
}

func Setup(c *RouteConfig) *gin.Engine {
    s := gin.New()
    s.ContextWithFallback = true

    s.Use(gin.Recovery())
    s.Use(middleware.LoggerMiddleware())

    api := s.Group("/api/v1")
    SetupAuthRoutes(api, c)
    return s
}

func SetupAuthRoutes(s *gin.RouterGroup, c *RouteConfig) {
	auth := s.Group("/auth")
	auth.GET("/", c.AuthHandler.GetUserDataHandler)
}
