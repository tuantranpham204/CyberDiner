package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/payload/request"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/payload/response"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/service"
)

type AuthController struct {
	auth service.AuthService
}

func NewAuthController(auth service.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

// SignOut revokes the current session's JWT by adding its jti to the denylist
// (UC-03). The Auth middleware must run before this handler.
func (ctl *AuthController) SignOut(c *gin.Context) {
	jti := middleware.TokenID(c)
	exp := middleware.TokenExpiresAt(c)

	if err := ctl.auth.SignOut(c.Request.Context(), jti, exp); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewError(
			http.StatusInternalServerError,
			"Unable to revoke session. Please try again later.",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK,
		"Signed out successfully",
		nil,
	))
}

func (ctl *AuthController) SignIn(c *gin.Context) {
	var req request.SignIn
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusUnprocessableEntity, response.NewError(
				http.StatusUnprocessableEntity,
				"One or more fields are invalid",
				validationErrorMessages(ve),
			))
			return
		}
		c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			err.Error(),
			nil,
		))
		return
	}

	result, err := ctl.auth.SignIn(c.Request.Context(), service.SignInInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, response.NewError(
				http.StatusUnauthorized,
				"Username or password is incorrect.",
				nil,
			))
		case errors.Is(err, service.ErrAccountBanned):
			c.JSON(http.StatusForbidden, response.NewError(
				http.StatusForbidden,
				"This account has been suspended. Please contact support.",
				nil,
			))
		case errors.Is(err, service.ErrAccountInactive):
			c.JSON(http.StatusForbidden, response.NewError(
				http.StatusForbidden,
				"This account is not active.",
				nil,
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewError(
				http.StatusInternalServerError,
				"Unable to sign in. Please try again later.",
				nil,
			))
		}
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK,
		"Signed in successfully",
		response.Auth{
			Token: result.Token,
			User:  response.FromUser(result.User),
		},
	))
}

func (ctl *AuthController) SignUp(c *gin.Context) {
	var req request.SignUp
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusUnprocessableEntity, response.NewError(
				http.StatusUnprocessableEntity,
				"One or more fields are invalid",
				validationErrorMessages(ve),
			))
			return
		}
		c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			err.Error(),
			nil,
		))
		return
	}

	result, err := ctl.auth.SignUp(c.Request.Context(), service.SignUpInput{
		Name:            req.Name,
		Surname:         req.Surname,
		Username:        req.Username,
		Email:           req.Email,
		DateOfBirth:     req.DateOfBirth,
		Gender:          req.Gender,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, response.NewError(
				http.StatusConflict,
				"An account with this email already exists. Please sign in instead.",
				nil,
			))
		case errors.Is(err, service.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, response.NewError(
				http.StatusConflict,
				"This username is already taken.",
				nil,
			))
		case errors.Is(err, service.ErrPasswordMismatch):
			c.JSON(http.StatusUnprocessableEntity, response.NewError(
				http.StatusUnprocessableEntity,
				"Password and confirmation do not match.",
				map[string]string{"confirm_password": "must match password"},
			))
		case errors.Is(err, service.ErrInvalidGender):
			c.JSON(http.StatusUnprocessableEntity, response.NewError(
				http.StatusUnprocessableEntity,
				"Gender value is not allowed.",
				map[string]string{"gender": "must be one of: male, female, other, prefer_not_to_say"},
			))
		case errors.Is(err, service.ErrInvalidDateOfBirth):
			c.JSON(http.StatusUnprocessableEntity, response.NewError(
				http.StatusUnprocessableEntity,
				"Date of birth is invalid.",
				map[string]string{"dob": "must be YYYY-MM-DD"},
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewError(
				http.StatusInternalServerError,
				"Unable to create account. Please try again later.",
				nil,
			))
		}
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccess(
		http.StatusCreated,
		"Account created successfully",
		response.Auth{
			Token: result.Token,
			User:  response.FromUser(result.User),
		},
	))
}

func validationErrorMessages(ve validator.ValidationErrors) map[string]string {
	out := make(map[string]string, len(ve))
	for _, fe := range ve {
		out[jsonFieldName(fe.Field())] = humanMessage(fe)
	}
	return out
}

func jsonFieldName(f string) string {
	switch f {
	case "Name":
		return "name"
	case "Surname":
		return "surname"
	case "Username":
		return "username"
	case "Email":
		return "email"
	case "DateOfBirth":
		return "dob"
	case "Gender":
		return "gender"
	case "Password":
		return "password"
	case "ConfirmPassword":
		return "confirm_password"
	case "PhoneNumber":
		return "phone_number"
	case "Address":
		return "address"
	}
	return f
}

func humanMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "username":
		return "must be 3-30 chars; letters, digits, dot, underscore, hyphen"
	case "strongpwd":
		return "must be 8-72 chars and include upper, lower, digit, and special character"
	case "eqfield":
		return "must match password"
	case "personname":
		return "must contain only letters, spaces, hyphen, or apostrophe"
	case "gender":
		return "must be one of: male, female, other, prefer_not_to_say"
	case "dob":
		return "must be a valid date (YYYY-MM-DD); age must be between 13 and 120"
	case "max":
		return "exceeds maximum length"
	case "phone":
		return "must be a valid Vietnamese phone number (e.g. 0912345678 or +84912345678)"
	}
	return "is invalid"
}
