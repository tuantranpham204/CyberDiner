package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/payload/request"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model/payload/response"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/service"
)

type ProfileController struct {
	profile service.ProfileService
}

func NewProfileController(profile service.ProfileService) *ProfileController {
	return &ProfileController{profile: profile}
}

// Get returns the profile of the user identified by the `:id` path parameter
// (UC-04 step 2). Requires authentication; any signed-in user may look up
// another user's profile.
func (ctl *ProfileController) Get(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, response.NewError(
			http.StatusBadRequest,
			"Path parameter `id` must be a positive integer.",
			nil,
		))
		return
	}

	user, err := ctl.profile.Get(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProfileNotFound):
			c.JSON(http.StatusNotFound, response.NewError(
				http.StatusNotFound, "Profile could not be located.", nil,
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewError(
				http.StatusInternalServerError, "Failed to load profile.", nil,
			))
		}
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK, "Profile retrieved",
		response.FromUser(user),
	))
}

// Update applies a partial profile update (UC-04 basic flow + Invalid Phone Number alt-flow).
func (ctl *ProfileController) Update(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.NewError(
			http.StatusUnauthorized, "Token does not identify a valid user.", nil,
		))
		return
	}

	var req request.UpdateProfile
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
			http.StatusBadRequest, err.Error(), nil,
		))
		return
	}

	user, err := ctl.profile.Update(c.Request.Context(), userID, service.UpdateProfileInput{
		Name:        req.Name,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNothingToUpdate):
			c.JSON(http.StatusBadRequest, response.NewError(
				http.StatusBadRequest,
				"Request body did not include any updatable fields.",
				nil,
			))
		case errors.Is(err, service.ErrProfileNotFound):
			c.JSON(http.StatusNotFound, response.NewError(
				http.StatusNotFound, "Profile could not be located.", nil,
			))
		default:
			c.JSON(http.StatusInternalServerError, response.NewError(
				http.StatusInternalServerError, "Failed to update profile.", nil,
			))
		}
		return
	}

	c.JSON(http.StatusOK, response.NewSuccess(
		http.StatusOK, "Profile updated successfully",
		response.FromUser(user),
	))
}
