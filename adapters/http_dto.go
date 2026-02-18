package adapters

// --- Auth DTOs ---

type AuthRequest struct {
	Email    string `json:"email" example:"test@example.com"`
	Password string `json:"password" example:"password123"`
}

type AuthResponse struct {
	Message string `json:"message" example:"operation successful"`
	Error   string `json:"error" example:""`
}

// --- Item DTOs ---

type CreateItemRequest struct {
	ProductName string `json:"productName" example:"iphone 71"`
	ProductDesc string `json:"productDesc" example:"Latest model with 128GB storage"`
}

type UpdateItemRequest struct {
	ProductName string `json:"productName" example:"iphone 71"`
	ProductDesc string `json:"productDesc" example:"Updated model with 256GB storage"`
}

type ItemResponse struct {
	ProductID   uint   `json:"productID" example:"1"`
	ProductName string `json:"productName" example:"iphone 71"`
	ProductDesc string `json:"productDesc" example:"Latest model"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"item not found"`
}
