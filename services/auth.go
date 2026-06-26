package services

import (
	"errors"

	"github.com/dotping-me/baastard/payload"
	"github.com/dotping-me/baastard/schema"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// This service is responsible for business logic regarding user
// authentication and authorisation.

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

// Method to register a new auth user onto the system:
func (s *AuthService) Register(req payload.RegisterAuthRequest) (*schema.Auth, error) {
	var existing schema.Auth

	// Checks if email already exists
	err := s.db.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		return nil, errors.New("Email already exists!")
	}

	// Any other errors that may arise regarding the database
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// TODO: Extra step to verify email address, if enabled

	// Hashes the password and creates the Auth user
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := schema.Auth{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	// Stores user in the database
	err = s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Method that allows a user to login in and creates a session
func (s *AuthService) Login(req payload.LoginRequest) (*schema.Auth, error) {
	var user schema.Auth

	err := s.db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, errors.New("Invalid credentials!")
	}

	// TODO: Add a check to verify is the user is not already logged in through
	// 		 querying sessions.

	// Validates Password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("Invalid credentials!")
	}

	return &user, nil
}

// Helper method that just queries for a user by email
func (s *AuthService) FindByEmail(email string) (*schema.Auth, error) {
	var user schema.Auth

	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Helper method that just queries for a user by id
func (s *AuthService) FindByID(id string) (*schema.Auth, error) {
	var user schema.Auth

	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
