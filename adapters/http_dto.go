package adapters

type RegisterRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"registered"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}
