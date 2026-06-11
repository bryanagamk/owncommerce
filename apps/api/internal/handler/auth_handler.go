package handler

import (
	"errors"

	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authSvc *auth.Service
}

func NewAuthHandler(authSvc *auth.Service) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	StoreName string `json:"store_name"`
	Slug      string `json:"slug"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Email == "" || req.Password == "" || req.Name == "" || req.StoreName == "" || req.Slug == "" {
		return response.BadRequest(c, "email, password, name, store_name, and slug are required")
	}

	result, err := h.authSvc.Register(c.UserContext(), auth.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		StoreName: req.StoreName,
		Slug:      req.Slug,
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
	})
	if err != nil {
		if errors.Is(err, auth.ErrEmailTaken) {
			return response.BadRequest(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, mapAuthResult(result))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "email and password are required")
	}

	result, err := h.authSvc.Login(c.UserContext(), auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return response.Unauthorized(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, mapAuthResult(result))
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.RefreshToken == "" {
		return response.BadRequest(c, "refresh_token is required")
	}

	result, err := h.authSvc.Refresh(c.UserContext(), auth.RefreshInput{
		RefreshToken: req.RefreshToken,
		UserAgent:    c.Get("User-Agent"),
		IPAddress:    c.IP(),
	})
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, mapAuthResult(result))
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	user := middleware.UserFromLocals(c)
	if user == nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req logoutRequest
	_ = c.BodyParser(&req)

	if err := h.authSvc.Logout(c.UserContext(), user.ID, req.RefreshToken, c.IP(), c.Get("User-Agent")); err != nil {
		return response.InternalError(c, "logout failed")
	}

	return response.Message(c, "logged out")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := middleware.UserFromLocals(c)
	if user == nil {
		return response.Unauthorized(c, "authentication required")
	}

	u, err := h.authSvc.GetUser(c.UserContext(), user.ID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	return response.OK(c, fiber.Map{
		"user": fiber.Map{
			"id":    u.ID,
			"email": u.Email,
			"name":  u.Name,
			"phone": u.Phone,
		},
		"tenant_id":   user.TenantID,
		"roles":       user.Roles,
		"permissions": user.Permissions,
	})
}

func mapAuthResult(result *auth.AuthResult) fiber.Map {
	data := fiber.Map{
		"user": fiber.Map{
			"id":    result.User.ID,
			"email": result.User.Email,
			"name":  result.User.Name,
		},
		"tokens":      result.Tokens,
		"roles":       result.Roles,
		"permissions": result.Permissions,
	}
	if result.Tenant != nil {
		data["tenant"] = fiber.Map{
			"id":   result.Tenant.ID,
			"name": result.Tenant.Name,
			"slug": result.Tenant.Slug,
		}
	}
	return data
}
