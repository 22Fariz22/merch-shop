package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/auth"
	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			mw.logger.Infof("auth middleware bearerHeader %s", bearerHeader)

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, mw.authUC, c, mw.cfg); err != nil {
					mw.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				return next(c)
			}

			cookie, err := c.Cookie("jwt-token")
			if err != nil {
				mw.logger.Errorf("c.Cookie", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			if err = mw.validateJWTToken(cookie.Value, mw.authUC, c, mw.cfg); err != nil {
				mw.logger.Errorf("validateJWTToken", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}
			return next(c)
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		fmt.Println("userUUID in validateJWTToken: ", userUUID)
		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		// req := c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
