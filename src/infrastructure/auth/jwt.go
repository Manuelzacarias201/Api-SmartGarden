package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService implementa la autenticación JWT
type JWTService struct {
	secretKey   string
	expiryHours int
}

// JWTClaims contiene los claims del token JWT
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWTService crea una nueva instancia de JWTService
func NewJWTService(secretKey string, expiryHours int) *JWTService {
	return &JWTService{
		secretKey:   secretKey,
		expiryHours: expiryHours,
	}
}

// GenerateToken genera un token JWT para un usuario
func (s *JWTService) GenerateToken(userID uint, username string) (string, error) {
	// Crear claims con datos de usuario
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear token con claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la clave secreta
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken valida un token JWT y devuelve el ID del usuario
func (s *JWTService) ValidateToken(tokenString string) (uint, error) {
	// Analizar el token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el algoritmo de firma es el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}

		return []byte(s.secretKey), nil
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
