package handler

import (
	"strconv"

	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
)

// ListSuppliers 获取供应商列表
func (h *Handler) ListSuppliers(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, "查询参数错误")
		return
	}
	query.GetOffset()

	suppliers, total, err := h.svc.ListSuppliers(&query)
	if err != nil {
		Error(c, 500, "获取供应商列表失败")
		return
	}
	Paginated(c, suppliers, total, query.Page, query.PageSize)
}

// GetAllSuppliers 获取所有启用供应商 (下拉选择用)
func (h *Handler) GetAllSuppliers(c *gin.Context) {
	suppliers, err := h.svc.GetAllSuppliers()
	if err != nil {
		Error(c, 500, "获取供应商失败")
		return
	}
	Success(c, suppliers)
}

// CreateSupplier 创建供应商
func (h *Handler) CreateSupplier(c *gin.Context) {
	var req models.SupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	sup, err := h.svc.CreateSupplier(&req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Created(c, sup)
}

// UpdateSupplier 更新供应商
func (h *Handler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的供应商ID")
		return
	}

	var req models.SupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	sup, err := h.svc.UpdateSupplier(uint(id), &req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, sup)
}

// DeleteSupplier 删除供应商
func (h *Handler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的供应商ID")
		return
	}

	if err := h.svc.DeleteSupplier(uint(id)); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, nil)
}
