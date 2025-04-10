package use_case

import (
	"context"
	"errors"
	"log"
	"time"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/events"
	"ApiSmart/src/core/domain/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthUseCase implementa los casos de uso relacionados con autenticación
type AuthUseCase struct {
	userRepo        application.UserRepository
	eventDispatcher application.EventDispatcher
	jwtService      JWTService
}

// JWTService define la interfaz para el servicio JWT
type JWTService interface {
	GenerateToken(userID uint, username string) (string, error)
	ValidateToken(token string) (uint, error)
}

// NewAuthUseCase crea una nueva instancia de AuthUseCase
func NewAuthUseCase(
	userRepo application.UserRepository,
	eventDispatcher application.EventDispatcher,
	jwtService JWTService,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:        userRepo,
		eventDispatcher: eventDispatcher,
		jwtService:      jwtService,
	}
}

// Register registra un nuevo usuario
func (uc *AuthUseCase) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Comprobar si el usuario ya existe
	existingUser, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("el correo electrónico ya está registrado")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Crear nuevo usuario
	now := time.Now()
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Guardar en base de datos
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generar token JWT
	token, err := uc.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// Publicar evento de registro
	if uc.eventDispatcher != nil {
		eventData := map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		}

		if err := uc.eventDispatcher.Dispatch(
			ctx,
			events.EventTypeUserRegistered,
			events.TopicUserEvents,
			eventData,
		); err != nil {
			log.Printf("Error al publicar evento de registro: %v", err)
		}
	}

	return &models.AuthResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// Login autentica a un usuario
func (uc *AuthUseCase) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	// Buscar usuario por email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := uc.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// Publicar evento de autenticación
	if uc.eventDispatcher != nil {
		eventData := map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		}

		if err := uc.eventDispatcher.Dispatch(
			ctx,
			events.EventTypeUserAuthenticated,
			events.TopicUserEvents,
			eventData,
		); err != nil {
			log.Printf("Error al publicar evento de autenticación: %v", err)
		}
	}

	return &models.AuthResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// ValidateToken valida un token JWT
func (uc *AuthUseCase) ValidateToken(token string) (uint, error) {
	return uc.jwtService.ValidateToken(token)
}
