package adapters

type RegisterRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"registered"`
}

type CreateItemRequest struct {
	Description string
}

type CreateItemResponse struct {
	Message string `json:"message" example:"registered"`
}

type ReadItemResponse struct {
	Message string `json:"message" example:"registered"`
}

type UpdateItemRequest struct {
	Description string
}

type UpdateItemResponse struct {
	Message string `json:"message" example:"registered"`
}

type DeleteItemResponse struct {
	Message string `json:"message" example:"registered"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}
