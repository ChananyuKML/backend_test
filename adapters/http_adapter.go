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
// @Param        request  body      AuthRequest  true  "Registration Info"
// @Success      200      {object}  map[string]string "message: registered"
// @Failure      400      {object}  map[string]string "message: registration failed"
// @Router       /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&req)

	err := h.uc.Register(req.Email, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "registration failed ",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "registered ",
		"error":   " ",
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and set auth_token and ref_token cookies
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      AuthRequest  true  "Login Credentials"
// @Success      200      {object}  map[string]string "message: login successfully"
// @Failure      401      {object}  map[string]string "message: fail to login"
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&req)

	access, refresh, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "fail to login",
			"error":   err.Error(),
		})
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
		"error":   " ",
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
		"error":   " ",
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
// @Description  Add a new item to the store
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request body      CreateItemRequest  true "Item Details"
// @Success      201     {string}  map[string]string "message: item created"
// @Failure      400     {object}  map[string]string "error: invalid request body"
// @Failure      500     {object}  map[string]string "error: failed to create item"
// @Router       /items [post]
func (h *ItemHandler) Create(c *fiber.Ctx) error {

	var req struct {
		ProductName string `json:"productName"`
		ProductDesc string `json:"productDesc"`
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
	return c.JSON(fiber.Map{
		"message": "item created",
		"error":   "",
	})
}

// List godoc
// @Summary      List all items
// @Description  Fetch all products from the database
// @Tags         items
// @Produce      json
// @Success      200  {object}  map[string]interface{} "message: [items...]"
// @Failure      500  {object}  map[string]interface{}
// @Router       /items [get]
func (h *ItemHandler) List(c *fiber.Ctx) error {
	items, err := h.uc.GetAllItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": []interface{}{},
			"error":   err,
		})
	}

	return c.JSON(fiber.Map{
		"message": items,
		"error":   "",
	})
}

// Update godoc
// @Summary      Update Item
// @Description  Update product name and description by ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id      path      int          true  "Product ID" example(1)
// @Param        request body      UpdateItemRequest  true  "New Item Data"
// @Success      200     {object}  map[string]string "message: item updated"
// @Failure      400     {object}  map[string]string "error: Invalid ID format"
// @Failure      403     {object}  map[string]string "error: forbidden"
// @Router       /items/{id} [put]
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req struct {
		ProductName string `json:"productName"`
		ProductDesc string `json:"productDesc"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": " ",
			"error":   err,
		})
	}

	if err := h.uc.UpdateItem(uint(id), req.ProductName, req.ProductDesc); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": " ",
			"error":   err,
		})
	}

	return c.JSON(fiber.Map{
		"message": "item updated",
		"error":   "",
	})
}

// Delete godoc
// @Summary      Delete an item
// @Description  Remove a product by ID
// @Tags         items
// @Param        id   path      int  true  "Item ID" example(1)
// @Success      200  {object}  map[string]string "message: item deleted"
// @Failure      400  {object}  map[string]string "error: Invalid ID format"
// @Failure      403  {object}  map[string]string "error: forbidden"
// @Router       /items/{id} [delete]
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": " ",
			"error":   "Invalid ID format",
		})
	}

	if err := h.uc.DeleteItem(uint(id)); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": " ",
			"error":   err,
		})
	}

	return c.JSON(fiber.Map{
		"message": "item deleted",
		"error":   "",
	})
}
