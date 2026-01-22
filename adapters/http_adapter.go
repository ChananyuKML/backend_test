package adapters

import (
	"hole/use_cases"
	"strconv"

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
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register payload"
// @Success 200 {object} adapters.RegisterResponse "Registration successful"
// @Failure 401 {object} adapters.ErrorResponse "Invalid credentials"
// @Router /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&req)

	if err := h.uc.Register(req.Email, req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "registered"})
}

// Login godoc
// @Summary Login user
// @Description Login and get access & refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} adapters.LoginResponse "Login successful"
// @Failure 401 {object} adapters.ErrorResponse "Invalid credentials"
// @Router /login [post]
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

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	access, refresh, err := h.uc.Refresh(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *ItemHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req struct {
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "description is required",
		})
	}

	if err := h.uc.CreateItem(userID, req.Description); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create item",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *ItemHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	items, err := h.uc.GetMyItems(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch items",
		})
	}

	return c.JSON(items)
}

func (h *ItemHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	idParam := c.Params("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid item id",
		})
	}

	var req struct {
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.uc.UpdateItem(uint(itemID), userID, req.Description); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "item not found or not owned by user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	idParam := c.Params("id")
	itemID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid item id",
		})
	}

	if err := h.uc.DeleteItem(uint(itemID), userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "item not found or not owned by user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
