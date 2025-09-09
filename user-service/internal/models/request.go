package models

import "github.com/auroban/gochat-be/user-service/internal/repository"

type ReqCreateUser struct {
	Username     string `json:"username" validate:"required,min=3,max=32"`
	Password     string `json:"password" validate:"required,min=8"`
	Email        string `json:"email" validate:"required,email"`
	DisplayName  string `json:"display_name" validate:"required"`
	AuthProvider string `json:"auth_provider" validate:"omitempty,oneof=google github local"`
	Type         string `json:"type" validate:"required,oneof=admin user"`
	Status       string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
}

func (r *ReqCreateUser) ToUser() *repository.User {
	return &repository.User{
		Username:     r.Username,
		PasswordHash: "",
		Email:        r.Email,
		DisplayName:  r.DisplayName,
		AuthProvider: r.AuthProvider,
		Type:         r.Type,
		Status:       r.Status,
	}
}
