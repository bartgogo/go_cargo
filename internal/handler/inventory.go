package handler

import (
	"strconv"

	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
)

// StockIn 入库操作
func (h *Handler) StockIn(c *gin.Context) {
	var req models.StockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := GetCurrentUserID(c)
	username := GetCurrentUsername(c)

	if err := h.svc.StockIn(&req, userID, username); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"message": "入库成功"})
}

// StockOut 出库操作
func (h *Handler) StockOut(c *gin.Context) {
	var req models.StockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := GetCurrentUserID(c)
	username := GetCurrentUsername(c)

	if err := h.svc.StockOut(&req, userID, username); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"message": "出库成功"})
}

// StockAdjust 库存调整
func (h *Handler) StockAdjust(c *gin.Context) {
	var req models.StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := GetCurrentUserID(c)
	username := GetCurrentUsername(c)

	if err := h.svc.StockAdjust(&req, userID, username); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"message": "库存调整成功"})
}

// ListInventoryRecords 查询库存操作记录
func (h *Handler) ListInventoryRecords(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		BadRequest(c, "查询参数错误")
		return
	}
	query.GetOffset()

	var productID *uint
	if pid := c.Query("product_id"); pid != "" {
		if id, err := strconv.ParseUint(pid, 10, 32); err == nil {
			uid := uint(id)
			productID = &uid
		}
	}

	recordType := c.Query("type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	records, total, err := h.svc.ListInventoryRecords(&query, productID, recordType, startDate, endDate)
	if err != nil {
		Error(c, 500, "获取库存记录失败")
		return
	}
	Paginated(c, records, total, query.Page, query.PageSize)
}
