package handler

import (
	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请输入用户名和密码")
		return
	}

	token, user, err := h.svc.Login(&req)
	if err != nil {
		Error(c, 401, err.Error())
		return
	}

	Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, user)
}

// GetProfile 获取个人信息
func (h *Handler) GetProfile(c *gin.Context) {
	userID := GetCurrentUserID(c)
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		Error(c, 404, "用户不存在")
		return
	}
	Success(c, user)
}

// UpdateProfile 更新个人信息
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}

	userID := GetCurrentUserID(c)
	user, err := h.svc.UpdateProfile(userID, &req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, user)
}

// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请输入原密码和新密码")
		return
	}

	userID := GetCurrentUserID(c)
	if err := h.svc.ChangePassword(userID, &req); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, nil)
}
