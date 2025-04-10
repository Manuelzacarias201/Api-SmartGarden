package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ApiSmart/src/core/application"
	"ApiSmart/src/core/domain/models"
)

// UserRepository implementa application.UserRepository
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) application.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create guarda un nuevo usuario en la base de datos
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

// FindByEmail busca un usuario por su correo electrónico
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at 
		FROM users 
		WHERE email = ?
	`

	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	var createdAt, updatedAt string

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Convertir strings a time.Time
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &user, nil
}

// FindByID busca un usuario por su ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	var createdAt, updatedAt string

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Convertir strings a time.Time
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &user, nil
}
