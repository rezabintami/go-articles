package middleware

import (
	"fmt"
	"go-articles/helpers"
	"net/http"
	"time"

	"go-articles/logger"

	"github.com/golang-jwt/jwt"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type JwtCustomClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type ConfigJWT struct {
	SecretJWT        string
	RefreshSecretJWT string
	ExpiredDuration  int
}

type ConfigForgotJWT struct {
	ResetSecretJWT  string
	ExpiredDuration int
}

type JwtRefreshToken struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type JwtResetPasswordToken struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func (config *ConfigJWT) Init() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(config.SecretJWT),
	}
}

func (config *ConfigForgotJWT) VerifyTokenForgotPassword() middleware.JWTConfig {
	return middleware.JWTConfig{
		ContextKey: "user_reset_password",
		Claims:     &JwtResetPasswordToken{},
		SigningKey: []byte(config.ResetSecretJWT),
	}
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if result, ok := err.(*echo.HTTPError); ok {
		helpers.ErrorResponse(c, result.Code, fmt.Errorf("%v", result.Message))
	} else {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err)
	}
}

// GenerateToken jwt ...
func (config *ConfigJWT) GenerateToken(userID int, name, role string) string {
	claims := JwtCustomClaims{
		userID,
		name,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(int64(config.ExpiredDuration))).Unix(),
		},
	}

	// Create token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := t.SignedString([]byte(config.SecretJWT))

	return token
}

// GenerateRefreshToken jwt ...
func (config *ConfigJWT) GenerateRefreshToken(userID int) string {
	claims := JwtRefreshToken{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(int64(config.ExpiredDuration))).Unix(),
		},
	}

	// Create token with claims
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := t.SignedString([]byte(config.RefreshSecretJWT))

	return token
}

func (config *ConfigForgotJWT) GenerateTokenResetPassword(id int, expHour time.Duration) (string, string, error) {
	expirationTime := time.Now().Add(expHour * time.Hour * 1).Unix()
	unixTimeUTC := time.Unix(expirationTime, 0)
	unitTimeInRFC3339 := unixTimeUTC.UTC().Format(time.RFC3339) // converts utc time to RFC3339 format

	claims := &JwtResetPasswordToken{
		id,
		jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(config.ResetSecretJWT))

	return token, unitTimeInRFC3339, err
}

// GetUser from jwt ...
func GetUser(c echo.Context) *JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims
}

// GetUserIdResetPassword from jwt ...
func GetUserIdResetPassword(c echo.Context) *JwtResetPasswordToken {
	user := c.Get("user_reset_password").(*jwt.Token)
	claims := user.Claims.(*JwtResetPasswordToken)
	return claims
}

func MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Logging(c)
		return next(c)
	}
}
