package server

import (
	authHttp "github.com/22Fariz22/merch-shop/internal/auth/delivery/http"
	authRepository "github.com/22Fariz22/merch-shop/internal/auth/repository"
	authUseCase "github.com/22Fariz22/merch-shop/internal/auth/usecase"

	apiMiddlewares "github.com/22Fariz22/merch-shop/internal/middleware"

	sessionRepository "github.com/22Fariz22/merch-shop/internal/session/repository"
	sessionUseCase "github.com/22Fariz22/merch-shop/internal/session/usecase"

	shopHttp "github.com/22Fariz22/merch-shop/internal/shop/delivery/http"
	shopRepository "github.com/22Fariz22/merch-shop/internal/shop/repository"
	shopUseCase "github.com/22Fariz22/merch-shop/internal/shop/usecase"

	transferHttp "github.com/22Fariz22/merch-shop/internal/transfer/delivery/http"
	transferRepository "github.com/22Fariz22/merch-shop/internal/transfer/repository"
	transferUseCase "github.com/22Fariz22/merch-shop/internal/transfer/usecase"

	infoHttp "github.com/22Fariz22/merch-shop/internal/info/delivery/http"
	infoRepository "github.com/22Fariz22/merch-shop/internal/info/repository"
	infoUseCase "github.com/22Fariz22/merch-shop/internal/info/usecase"

	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db, s.logger)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg, s.logger)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient, s.logger)
	shRepo := shopRepository.NewShopRepository(s.db, s.logger)
	trRepo := transferRepository.NewTransferRepository(s.db, s.logger)
	infoRepo := infoRepository.NewInfoRepository(s.db, s.logger)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	sessUC := sessionUseCase.NewSessionUseCase(sRepo, s.cfg)
	shopUC := shopUseCase.NewShopUseCase(s.cfg, shRepo, s.logger)
	transferUC := transferUseCase.NewTransferUseCase(s.cfg, trRepo, s.logger)
	infoUC := infoUseCase.NewInfoUseCase(s.cfg, infoRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	shopHandlers := shopHttp.NewShopHandlers(s.cfg, shopUC, s.logger)
	transferHandlers := transferHttp.NewTransferHandlers(s.cfg, transferUC, s.logger)
	infoHandlers := infoHttp.NewHandlers(s.cfg, infoUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	shopGroup := v1.Group("/buy")
	transferGroup := v1.Group("/sendCoin")
	infoGroup := v1.Group("/info")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	shopHttp.MapShopRoutes(shopGroup, shopHandlers, mw)
	transferHttp.MapTransferRoutes(transferGroup, transferHandlers, mw)
	infoHttp.MapInfoRoutes(infoGroup, infoHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
