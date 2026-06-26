package payload

// Expected JSON payload structure for registering an auth user
type RegisterAuthRequest struct {
	Email    string `json:"email" 		binding:"required,email"`
	Password string `json:"password" 	binding:"required,min=8"`
}
