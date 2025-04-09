package services

import (
	"context"
	"errors"
	"time"

	"ApiSmart/internal/core/domain"
	"ApiSmart/internal/core/ports"
	"ApiSmart/pkg/auth"
	"ApiSmart/pkg/event"
	"golang.org/x/crypto/bcrypt"
)

// EventDrivenAuthService implementa AuthService con arquitectura basada en eventos
type EventDrivenAuthService struct {
	userRepo        ports.UserRepository
	eventDispatcher *event.EventDispatcher
}

// NewEventDrivenAuthService crea un nuevo servicio de autenticación basado en eventos
func NewEventDrivenAuthService(userRepo ports.UserRepository, eventDispatcher *event.EventDispatcher) ports.AuthService {
	return &EventDrivenAuthService{
		userRepo:        userRepo,
		eventDispatcher: eventDispatcher,
	}
}

// Register registra un nuevo usuario y publica un evento
func (s *EventDrivenAuthService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Comprobar si el usuario ya existe
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
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
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Guardar en base de datos
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generar token JWT
	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// Publicar evento de registro de usuario
	eventData := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	if err := s.eventDispatcher.Dispatch(
		ctx,
		event.EventTypeUserRegistered,
		event.TopicUserEvents,
		eventData,
	); err != nil {
		// Solo logear el error, no fallar el flujo principal
		// log.Printf("Error publishing user registered event: %v", err)
	}

	return &domain.AuthResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// Login autentica a un usuario y publica un evento
func (s *EventDrivenAuthService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Publicar evento de autenticación de usuario
	eventData := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}

	if err := s.eventDispatcher.Dispatch(
		ctx,
		event.EventTypeUserAuthenticated,
		event.TopicUserEvents,
		eventData,
	); err != nil {
		// Solo logear el error, no fallar el flujo principal
		// log.Printf("Error publishing user authenticated event: %v", err)
	}

	return &domain.AuthResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// ValidateToken valida un token JWT
func (s *EventDrivenAuthService) ValidateToken(token string) (uint, error) {
	return auth.ValidateJWT(token)
}
