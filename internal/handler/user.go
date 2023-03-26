package handler

import (
	"net/http"
	"server/internal/domain"
	"server/internal/port"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	port.UserServicePort
}

func NewUserHandler(s port.UserServicePort) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var u domain.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.UserServicePort.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Login(c *gin.Context) {
	var user domain.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.UserServicePort.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	res := &domain.LoginUserRes{
		AccessToken: u.AccessToken,
		ID:          u.ID,
		Username:    u.Username,
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var u domain.UpdateUsernameReq

	userID := c.Param("userId")
	u.ID = userID
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UserServicePort.UpdateUsername(c.Request.Context(), &u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "username updated successfully"})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserServicePort.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
