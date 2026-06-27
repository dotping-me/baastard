package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dotping-me/baastard/schema"
	"gorm.io/gorm"
)

// This service is responsible for managing sessions and access to the system.
// It is done through the use of Access Tokens and Refresh Tokens.

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

// Helper function to generate an random byte string to act as the refresh token.
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// Method that creates a session for a user
func (s *SessionService) Create(userID string,
	deviceID string, userAgent string, ip string) (*schema.Session, error) {

	// Generates a refresh token
	token, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create and store session in the database
	session := schema.Session{
		UserID:       userID,
		RefreshToken: token,

		DeviceID:  deviceID,
		UserAgent: userAgent,
		IPAddress: ip,

		ExpiresAt: time.Now().Add(RefreshCookieMaxAge),
	}

	err = s.db.Create(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Method to query a valid session by refresh token
func (s *SessionService) Get(token string) (*schema.Session, error) {

	var session schema.Session
	err := s.db.Where("refresh_token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	// Checks if token has expired
	if time.Now().After(session.ExpiresAt) {
		_ = s.db.Delete(&session).Error // Delete that expired session
		return nil, errors.New("Session expired!")
	}

	return &session, nil
}

// Method that returns all active sessions for a user
func (s *SessionService) GetAll(userID string) ([]schema.Session, error) {
	var sessions []schema.Session

	err := s.db.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// Method to issue a new refresh token for a valid session and invalidates
// the old refresh token
func (s *SessionService) Refresh(session *schema.Session) (*schema.Session, error) {
	newToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Updates the session
	session.RefreshToken = newToken
	session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

	// Updates record in database
	err = s.db.Save(session).Error
	if err != nil {
		return nil, err // TODO: Maybe add a custom error message here
	}

	return session, nil
}

// Method to delete a valid session
func (s *SessionService) Delete(token string) error {
	return s.db.Where("refresh_token=?", token).Delete(&schema.Session{}).Error
}

// Method to delete all sessions belonging to a user ID
func (s *SessionService) DeleteAll(userID string) error {
	return s.db.Where("user_id=?", userID).Delete(&schema.Session{}).Error
}

// Helper method to validate a device connection
func (s *SessionService) ValidateDevice(session *schema.Session, deviceID string) error {
	if session.DeviceID != deviceID {
		return errors.New("Invalid device!")
	}

	return nil
}
