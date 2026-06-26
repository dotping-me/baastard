package payload

// Expected JSON payload structure for logging in a user
type LoginRequest struct {
	Email    string `json:"email" 		binding:"required,email"`
	Password string `json:"password"	 binding:"required"`
}
