package payload

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
