package controllers

import (
	"kuliah-web-bsm-go/app/facades"
	"kuliah-web-bsm-go/app/models"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Login(ctx http.Context) http.Response {
	var req struct {
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid request parameters",
		})
	}

	if req.Email == "" || req.Password == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Email and password are required",
		})
	}

	var user models.User
	err := facades.Orm().Query().Where("email", req.Email).First(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Invalid credentials",
		})
	}

	if user.Password == "" || !facades.Hash().Check(req.Password, user.Password) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"token": token,
	})
}

func (c *AuthController) Register(ctx http.Context) http.Response {
	var req struct {
		Username string `form:"username" json:"username"`
		Email    string `form:"email" json:"email"`
		Password string `form:"password" json:"password"`
		Nama     string `form:"nama" json:"nama"`
		Telepon  string `form:"telepon" json:"telepon"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid request parameters",
		})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.Nama == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Username, email, password, and nama are required",
		})
	}

	// Check if email already exists
	var existing models.User
	err := facades.Orm().Query().Where("email", req.Email).First(&existing)
	if err == nil && existing.ID > 0 {
		return ctx.Response().Json(http.StatusConflict, http.Json{
			"message": "Email is already registered",
		})
	}

	hashedPassword, err := facades.Hash().Make(req.Password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to hash password",
		})
	}

	// Start database transaction
	user := models.User{
		RoleID:   2, // default customer
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = facades.Orm().Transaction(func(tx orm.Query) error {
		if err := tx.Create(&user); err != nil {
			return err
		}

		customer := models.Customer{
			UserID:  user.ID,
			Nama:    req.Nama,
			Telepon: &req.Telepon,
		}

		if err := tx.Create(&customer); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Registration failed: " + err.Error(),
		})
	}

	// Generate JWT token
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) GoogleCallback(ctx http.Context) http.Response {
	var req struct {
		GoogleID string `json:"google_id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Invalid request parameters",
		})
	}

	if req.GoogleID == "" || req.Email == "" || req.Username == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "google_id, email, and username are required",
		})
	}

	var user models.User
	err := facades.Orm().Query().Where("email", req.Email).First(&user)
	if err != nil {
		// User does not exist, let's register them
		// User does not exist, let's register them
		user = models.User{
			RoleID:   2, // default customer
			Username: req.Username,
			Email:    req.Email,
			GoogleID: &req.GoogleID,
			Avatar:   &req.Avatar,
		}

		err = facades.Orm().Transaction(func(tx orm.Query) error {
			if err := tx.Create(&user); err != nil {
				return err
			}

			customer := models.Customer{
				UserID: user.ID,
				Nama:   req.Username,
			}

			if err := tx.Create(&customer); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"message": "Google registration failed: " + err.Error(),
			})
		}
	} else {
		// User exists, update google_id and avatar if empty
		if user.GoogleID == nil || *user.GoogleID == "" {
			user.GoogleID = &req.GoogleID
			user.Avatar = &req.Avatar
			facades.Orm().Query().Save(&user)
		}
	}

	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) Profile(ctx http.Context) http.Response {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Unauthorized",
		})
	}

	// Load relationships
	err = facades.Orm().Query().With("Role").With("Customer").Where("id", user.ID).First(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "User not found",
		})
	}

	return ctx.Response().Success().Json(user)
}
