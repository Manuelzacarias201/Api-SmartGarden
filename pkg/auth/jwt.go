package auth

import (
	"errors"
	"time"

	"ApiSmart/config"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username string) (string, error) {
	// Cargar la configuración
	cfg := config.LoadConfig()

	// Crear claims con datos de usuario
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token válido por 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear token con claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (uint, error) {
	// Cargar la configuración
	cfg := config.LoadConfig()

	// Analizar el token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el algoritmo de firma es el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}

		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// Verificar si el token es válido
	if !token.Valid {
		return 0, errors.New("token inválido")
	}

	// Extraer claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, errors.New("no se pudieron extraer los claims")
	}

	return claims.UserID, nil
}
