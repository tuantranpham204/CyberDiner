package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) SignUp(c *gin.Context) {
	var req service.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := ctrl.authService.Register(&req); err != nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

func (ctrl *AuthController) SignIn(c *gin.Context) {
	var req service.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	resp, err := ctrl.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrAccountSuspended) {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// SignOut delegates to the service with the JTI and expiry that JWTAuth middleware
// already validated and placed in context — no JWT logic here.
func (ctrl *AuthController) SignOut(c *gin.Context) {
	jti, expiry := middleware.GetJTI(c)
	ctrl.authService.SignOut(jti, expiry)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "signed out"})
}
