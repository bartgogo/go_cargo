package handler

import (
	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计数据
func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		Error(c, 500, "获取统计数据失败")
		return
	}
	Success(c, stats)
}

// GetChartData 获取图表数据
func (h *Handler) GetChartData(c *gin.Context) {
	data, err := h.svc.GetChartData()
	if err != nil {
		Error(c, 500, "获取图表数据失败")
		return
	}
	Success(c, data)
}

// GetLowStockProducts 获取低库存预警
func (h *Handler) GetLowStockProducts(c *gin.Context) {
	products, err := h.svc.GetLowStockProducts()
	if err != nil {
		Error(c, 500, "获取低库存数据失败")
		return
	}
	Success(c, products)
}
