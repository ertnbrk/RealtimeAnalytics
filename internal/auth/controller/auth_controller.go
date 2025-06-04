package controller

import (
	"net/http"

	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/model"
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authSvc service.AuthService) *AuthController {
	return &AuthController{authService: authSvc}
}

func (ctr *AuthController) SignUpHandler(c *gin.Context) {
	var req model.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctr.authService.SignUp(&req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func (ctr *AuthController) LoginHandler(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctr.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kimlik bilgileri"})
		return
	}

	c.JSON(http.StatusOK, model.JWTResponse{Token: token})
}
