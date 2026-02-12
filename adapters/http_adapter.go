package adapters

import (
	"hole/use_cases"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	uc *use_cases.AuthUseCase
}

type ItemHandler struct {
	uc *use_cases.ItemUseCase
}

func NewAuthHandler(uc *use_cases.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc}
}

func NewItemHandler(uc *use_cases.ItemUseCase) *ItemHandler {
	return &ItemHandler{uc}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Registration Info"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Router       /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&req)

	err := h.uc.Register(req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "registered"})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and set auth_token and ref_token cookies
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Login Credentials"
// @Success      200      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&req)

	access, refresh, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	acc := new(fiber.Cookie)
	acc.Name = "auth_token"
	acc.Value = access
	acc.Expires = time.Now().Add(30 * time.Minute)
	acc.HTTPOnly = true
	acc.Secure = true
	acc.SameSite = "Lax"

	ref := new(fiber.Cookie)
	ref.Name = "ref_token"
	ref.Value = refresh
	ref.HTTPOnly = true
	ref.Secure = true
	ref.SameSite = "Lax"

	c.Cookie(acc)
	c.Cookie(ref)

	return c.JSON(fiber.Map{
		"message": "login sucessfully ",
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "ref_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

// func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
// 	var body struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}

// 	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid request body",
// 		})
// 	}

// 	access, refresh, err := h.uc.Refresh(body.RefreshToken)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"access_token":  access,
// 		"refresh_token": refresh,
// 	})
// }

// Create godoc
// @Summary      Create Item
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request body adapters.CreateItemRequest true "Item Details"
// @Success      201 {object} adapters.AuthResponse
// @Failure      400 {object} adapters.ErrorResponse
// @Router       /items [post]
func (h *ItemHandler) Create(c *fiber.Ctx) error {

	var req struct {
		ProductName string `json:"name"`
		ProductDesc string `json:"desc"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.ProductName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required is required",
		})
	}
	if req.ProductDesc == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product description is required",
		})
	}

	if err := h.uc.CreateItem(req.ProductName, req.ProductDesc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create item",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ListItems godoc
// @Summary      List all items
// @Description  Fetch all products from the database
// @Tags         items
// @Produce      json
// @Success      200  {array}   entities.Item
// @Failure      500  {object}  map[string]string
// @Router       /items [get]
func (h *ItemHandler) List(c *fiber.Ctx) error {

	items, err := h.uc.GetAllItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch items",
		})
	}

	return c.JSON(items)
}

// func (h *ItemHandler) List(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uint)

// 	items, err := h.uc.GetMyItems(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "failed to fetch items",
// 		})
// 	}

// 	return c.JSON(items)
// }

// Update godoc
// @Summary      Update Item
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id      path int                       true "Product ID" example(1)
// @Param        request body adapters.UpdateItemRequest true "New Item Data"
// @Success      200 {object} adapters.AuthResponse
// @Failure      404 {object} adapters.ErrorResponse
// @Router       /items/{id} [put]
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// idParam := c.Params("id")
	// itemID, err := strconv.Atoi(idParam)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "invalid item id",
	// 	})
	// }

	var req struct {
		ProductName string `json:"name"`
		ProductDesc string `json:"desc"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.uc.UpdateItem(uint(id), req.ProductName, req.ProductDesc); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "item not found",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteItem godoc
// @Summary      Delete an item
// @Description  Remove a product by ID
// @Tags         items
// @Param        id   path      int  true  "Item ID"
// @Success      204  {string}  string "No Content"
// @Failure      404  {object}  map[string]string
// @Router       /items/{id} [delete]
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	if err := h.uc.DeleteItem(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
