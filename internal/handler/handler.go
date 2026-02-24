// Package handler 提供通用 HTTP 处理工具
package handler

import (
	"math"
	"net/http"

	"go-cargo/internal/models"
	"go-cargo/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler 处理器基础结构
type Handler struct {
	svc *service.Service
}

// New 创建 Handler 实例
func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "created",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, models.Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 请求错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Paginated 分页响应
func Paginated(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	Success(c, models.PaginatedData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		return id.(uint)
	}
	return 0
}

// GetCurrentUsername 从上下文获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	if name, exists := c.Get("username"); exists {
		return name.(string)
	}
	return ""
}
