package handler

import (
	"strconv"

	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
)

// ListProducts 获取商品列表
func (h *Handler) ListProducts(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, "查询参数错误")
		return
	}
	query.GetOffset() // 初始化默认值

	var categoryID, supplierID *uint
	if cid := c.Query("category_id"); cid != "" {
		if id, err := strconv.ParseUint(cid, 10, 32); err == nil {
			uid := uint(id)
			categoryID = &uid
		}
	}
	if sid := c.Query("supplier_id"); sid != "" {
		if id, err := strconv.ParseUint(sid, 10, 32); err == nil {
			uid := uint(id)
			supplierID = &uid
		}
	}

	products, total, err := h.svc.ListProducts(&query, categoryID, supplierID)
	if err != nil {
		Error(c, 500, "获取商品列表失败")
		return
	}

	Paginated(c, products, total, query.Page, query.PageSize)
}

// GetProduct 获取商品详情
func (h *Handler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	product, err := h.svc.GetProduct(uint(id))
	if err != nil {
		Error(c, 404, "商品不存在")
		return
	}
	Success(c, product)
}

// CreateProduct 创建商品
func (h *Handler) CreateProduct(c *gin.Context) {
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, err := h.svc.CreateProduct(&req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Created(c, product)
}

// UpdateProduct 更新商品
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, err := h.svc.UpdateProduct(uint(id), &req)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, product)
}

// DeleteProduct 删除商品
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		BadRequest(c, "无效的商品ID")
		return
	}

	if err := h.svc.DeleteProduct(uint(id)); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, nil)
}
