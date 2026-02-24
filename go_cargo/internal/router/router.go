// Package router 定义 HTTP 路由
package router

import (
	"io/fs"
	"net/http"

	"go-cargo/internal/handler"
	"go-cargo/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 配置路由
func Setup(h *handler.Handler, webFS fs.FS) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 静态文件服务 (前端页面)
	r.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(webFS, "index.html")
		if err != nil {
			c.String(500, "页面加载失败")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
	r.GET("/app", func(c *gin.Context) {
		data, err := fs.ReadFile(webFS, "app.html")
		if err != nil {
			c.String(500, "页面加载失败")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由 (无需认证)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 个人信息
			protected.GET("/auth/profile", h.GetProfile)
			protected.PUT("/auth/profile", h.UpdateProfile)
			protected.PUT("/auth/change-password", h.ChangePassword)

			// 仪表盘
			protected.GET("/dashboard/stats", h.GetDashboardStats)
			protected.GET("/dashboard/charts", h.GetChartData)
			protected.GET("/dashboard/low-stock", h.GetLowStockProducts)

			// 分类管理
			protected.GET("/categories", h.ListCategories)
			protected.GET("/categories/all", h.GetAllCategories)
			protected.POST("/categories", h.CreateCategory)
			protected.PUT("/categories/:id", h.UpdateCategory)
			protected.DELETE("/categories/:id", h.DeleteCategory)

			// 供应商管理
			protected.GET("/suppliers", h.ListSuppliers)
			protected.GET("/suppliers/all", h.GetAllSuppliers)
			protected.POST("/suppliers", h.CreateSupplier)
			protected.PUT("/suppliers/:id", h.UpdateSupplier)
			protected.DELETE("/suppliers/:id", h.DeleteSupplier)

			// 商品管理
			protected.GET("/products", h.ListProducts)
			protected.GET("/products/:id", h.GetProduct)
			protected.POST("/products", h.CreateProduct)
			protected.PUT("/products/:id", h.UpdateProduct)
			protected.DELETE("/products/:id", h.DeleteProduct)

			// 库存操作
			protected.POST("/inventory/stock-in", h.StockIn)
			protected.POST("/inventory/stock-out", h.StockOut)
			protected.POST("/inventory/adjust", h.StockAdjust)
			protected.GET("/inventory/records", h.ListInventoryRecords)
		}
	}

	return r
}
