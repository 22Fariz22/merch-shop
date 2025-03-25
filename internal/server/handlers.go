package server

import (
	authHttp "github.com/22Fariz22/merch-shop/internal/auth/delivery/http"
	authRepository "github.com/22Fariz22/merch-shop/internal/auth/repository"
	authUseCase "github.com/22Fariz22/merch-shop/internal/auth/usecase"
	apiMiddlewares "github.com/22Fariz22/merch-shop/internal/middleware"
	sessionRepository "github.com/22Fariz22/merch-shop/internal/session/repository"
	"github.com/22Fariz22/merch-shop/internal/session/usecase"
	"github.com/22Fariz22/merch-shop/pkg/utils"
		"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)

	// Init useCases

	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
