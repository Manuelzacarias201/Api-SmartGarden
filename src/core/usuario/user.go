package usuario

import (
	"context"
	"time"
)

// UserContext representa el contexto del usuario en la aplicación
type UserContext struct {
	UserID      uint
	Username    string
	Email       string
	Roles       []string
	Permissions []string
	CreatedAt   time.Time
}

// WithUserContext agrega el contexto de usuario al contexto de Go
func WithUserContext(ctx context.Context, userCtx UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, userCtx)
}

// FromContext extrae el contexto de usuario del contexto de Go
func FromContext(ctx context.Context) (UserContext, bool) {
	userCtx, ok := ctx.Value(userContextKey).(UserContext)
	return userCtx, ok
}

// contextKey es un tipo para las claves de contexto
type contextKey string

// userContextKey es la clave para el contexto de usuario
const userContextKey = contextKey("user-context")
