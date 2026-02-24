package handler

import (
	"strconv"

	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
)

// ListCategories 获取分类列表
func (h *Handler) ListCategories(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, "查询参数错误")
		return
	}
	query.GetOffset()

	categories, total, err := h.svc.ListCategories(&query)
	if err != nil {
		Error(c, 500, "获取分类列表失败")
		return
	}
	Paginated(c, categories, total, query.Page, query.PageSize)
}

// GetAllCategories 获取所有启用分类 (下拉选择用)
func (h *Handler) GetAllCategories(c *gin.Context) {
	categories, err := h.svc.GetAllCategories()
	if err != nil {
		Error(c, 500, "获取分类失败")
		return
	}
	Success(c, categories)
}

// CreateCategory 创建分类
func (h *Handler) CreateCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	cat, err := h.svc.CreateCategory(&req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Created(c, cat)
}

// UpdateCategory 更新分类
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的分类ID")
		return
	}

	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	cat, err := h.svc.UpdateCategory(uint(id), &req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, cat)
}

// DeleteCategory 删除分类
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的分类ID")
		return
	}

	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, nil)
}
