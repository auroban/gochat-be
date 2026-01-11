package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/auroban/gochat-be/user-service/internal/models"
	"github.com/auroban/gochat-be/user-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	validator *validator.Validate
	service   *service.UserService
}

func NewUserHandler(validator *validator.Validate, service *service.UserService) *UserHandler {
	return &UserHandler{
		validator: validator,
		service:   service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.ReqCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateNewUser(r.Context(), req)
	if err != nil {
		logger.WithError(err).Error("error creating user")
		var httpStatus int
		var message string
		if errors.Is(err, service.ErrUserAlreadyExists) {
			httpStatus = http.StatusConflict
			message = err.Error()
		} else {
			httpStatus = http.StatusInternalServerError
			message = "Something went wrong"
		}
		http.Error(w, message, httpStatus)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}
	logger.Debugf("Retrieving user by ID: [%s]", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		logger.WithError(err).Errorf("error while fetching user by id: [%d]", id)
		var status int
		var msg string
		if errors.Is(err, service.ErrNoUserFound) {
			status = http.StatusNotFound
			msg = err.Error()
		} else {
			status = http.StatusInternalServerError
			msg = "Something went wrong"
		}
		http.Error(w, msg, status)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users", h.GetUserByID).Methods("GET")
}
