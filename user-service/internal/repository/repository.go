package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/auroban/gochat-be/user-service/internal/domainerror"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const QUERY_FIND_BY_ID = `SELECT * FROM "user" WHERE "id" = :id`
const QUERY_FIND_BY_USERNAME = `SELECT * FROM "user" WHERE username = :username`
const QUERY_FIND_BY_EMAIL = `SELECT * FROM "user" WHERE email = :email`
const QUERY_CREATE_USER = `
	INSERT INTO "user"
	(created_at, updated_at, created_by, updated_by, username, password_hash, email, display_name, auth_provider, type, status) 
	VALUES
	(:created_at, :updated_at, :created_by, :updated_by, :username, :password_hash, :email, :display_name, :auth_provider, :type, :status)
	RETURNING id
`

type BaseEntity struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy *int      `db:"created_by" json:"created_by,omitempty"` // nullable
	UpdatedBy *int      `db:"updated_by" json:"updated_by,omitempty"` // nullable
}

type User struct {
	BaseEntity
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"password_hash,omitempty"`
	Email        string `db:"email" json:"email"`
	DisplayName  string `db:"display_name" json:"display_name"`
	AuthProvider string `db:"auth_provider" json:"auth_provider,omitempty"`
	Type         string `db:"type" json:"type"`
	Status       string `db:"status" json:"status"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindById(ctx context.Context, id int) (*User, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, QUERY_FIND_BY_ID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.Get(&user, map[string]interface{}{"id": id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerror.ErrNoUserFound
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, QUERY_FIND_BY_USERNAME)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.Get(&user, map[string]interface{}{"username": username})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerror.ErrNoUserFound
		} else {
			return nil, err
		}
	}
	return &user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, QUERY_FIND_BY_EMAIL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user User
	err = stmt.Get(&user, map[string]interface{}{"email": email})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerror.ErrNoUserFound
		} else {
			return nil, err
		}
	}
	return &user, err
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	stmt, err := r.db.PrepareNamedContext(ctx, QUERY_CREATE_USER)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.Get(&id, user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("%w: constraint=%s detail=%s",
				domainerror.ErrUserAlreadyExists,
				pqErr.Constraint,
				pqErr.Detail,
			)
		}
		return nil, err
	}
	user.ID = id
	return user, nil
}
