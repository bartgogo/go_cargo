// Package models 定义系统数据模型
package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------- 基础模型 ----------

// BaseModel 通用基础模型，所有实体继承
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ---------- 用户模型 ----------

// User 用户
type User struct {
	BaseModel
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email    string `json:"email" gorm:"size:100"`
	Password string `json:"-" gorm:"size:255;not null"` // JSON 序列化时忽略密码
	RealName string `json:"real_name" gorm:"size:50"`
	Phone    string `json:"phone" gorm:"size:20"`
	Role     string `json:"role" gorm:"size:20;default:operator"` // admin, operator
	Status   int    `json:"status" gorm:"default:1"`              // 1=启用, 0=禁用
	Avatar   string `json:"avatar" gorm:"size:255"`
}

// TableName 指定表名
func (User) TableName() string { return "users" }

// ---------- 分类模型 ----------

// Category 商品分类
type Category struct {
	BaseModel
	Name        string `json:"name" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"size:500"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"` // 1=启用, 0=禁用

	// 关联统计 (不存储在数据库)
	ProductCount int64 `json:"product_count" gorm:"-"`
}

// TableName 指定表名
func (Category) TableName() string { return "categories" }

// ---------- 供应商模型 ----------

// Supplier 供应商
type Supplier struct {
	BaseModel
	Code          string `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name          string `json:"name" gorm:"size:200;not null"`
	ContactPerson string `json:"contact_person" gorm:"size:50"`
	Phone         string `json:"phone" gorm:"size:20"`
	Email         string `json:"email" gorm:"size:100"`
	Address       string `json:"address" gorm:"size:500"`
	Status        int    `json:"status" gorm:"default:1"` // 1=启用, 0=禁用
	Remark        string `json:"remark" gorm:"size:500"`

	// 关联统计
	ProductCount int64 `json:"product_count" gorm:"-"`
}

// TableName 指定表名
func (Supplier) TableName() string { return "suppliers" }

// ---------- 商品模型 ----------

// Product 商品
type Product struct {
	BaseModel
	SKU          string  `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Name         string  `json:"name" gorm:"size:200;not null;index"`
	Description  string  `json:"description" gorm:"size:1000"`
	CategoryID   *uint   `json:"category_id" gorm:"index"`
	SupplierID   *uint   `json:"supplier_id" gorm:"index"`
	Unit         string  `json:"unit" gorm:"size:20;default:个"` // 计量单位
	CostPrice    float64 `json:"cost_price" gorm:"type:decimal(12,2);default:0"`
	SellingPrice float64 `json:"selling_price" gorm:"type:decimal(12,2);default:0"`
	CurrentStock int     `json:"current_stock" gorm:"default:0"`
	MinStock     int     `json:"min_stock" gorm:"default:0"` // 最低库存预警
	MaxStock     int     `json:"max_stock" gorm:"default:0"` // 最高库存上限
	Barcode      string  `json:"barcode" gorm:"size:100;index"`
	Location     string  `json:"location" gorm:"size:100"` // 库位
	ImageURL     string  `json:"image_url" gorm:"size:500"`
	Status       int     `json:"status" gorm:"default:1"` // 1=启用, 0=禁用

	// 关联
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Supplier *Supplier `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
}

// TableName 指定表名
func (Product) TableName() string { return "products" }

// ---------- 库存记录模型 ----------

// InventoryRecordType 库存操作类型
type InventoryRecordType string

const (
	StockIn     InventoryRecordType = "stock_in"  // 入库
	StockOut    InventoryRecordType = "stock_out" // 出库
	StockAdjust InventoryRecordType = "adjust"    // 调整
)

// InventoryRecord 库存操作记录
type InventoryRecord struct {
	ID           uint                `json:"id" gorm:"primaryKey"`
	ProductID    uint                `json:"product_id" gorm:"index;not null"`
	Type         InventoryRecordType `json:"type" gorm:"size:20;not null;index"`
	Quantity     int                 `json:"quantity" gorm:"not null"` // 操作数量 (正数)
	BeforeQty    int                 `json:"before_qty"`               // 操作前数量
	AfterQty     int                 `json:"after_qty"`                // 操作后数量
	UnitCost     float64             `json:"unit_cost" gorm:"type:decimal(12,2);default:0"`
	TotalCost    float64             `json:"total_cost" gorm:"type:decimal(12,2);default:0"`
	ReferenceNo  string              `json:"reference_no" gorm:"size:100;index"` // 关联单号
	Notes        string              `json:"notes" gorm:"size:500"`
	OperatorID   uint                `json:"operator_id" gorm:"index"`
	OperatorName string              `json:"operator_name" gorm:"size:50"`
	CreatedAt    time.Time           `json:"created_at" gorm:"index"`

	// 关联
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName 指定表名
func (InventoryRecord) TableName() string { return "inventory_records" }

// ---------- API 请求/响应结构体 ----------

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email"`
	RealName string `json:"real_name"`
}

// UpdateProfileRequest 更新个人信息
type UpdateProfileRequest struct {
	Email    string `json:"email"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

// ChangePasswordRequest 修改密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// ProductRequest 商品请求
type ProductRequest struct {
	SKU          string  `json:"sku" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	CategoryID   *uint   `json:"category_id"`
	SupplierID   *uint   `json:"supplier_id"`
	Unit         string  `json:"unit"`
	CostPrice    float64 `json:"cost_price"`
	SellingPrice float64 `json:"selling_price"`
	MinStock     int     `json:"min_stock"`
	MaxStock     int     `json:"max_stock"`
	Barcode      string  `json:"barcode"`
	Location     string  `json:"location"`
	ImageURL     string  `json:"image_url"`
	Status       int     `json:"status"`
}

// CategoryRequest 分类请求
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// SupplierRequest 供应商请求
type SupplierRequest struct {
	Code          string `json:"code" binding:"required"`
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contact_person"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Status        int    `json:"status"`
	Remark        string `json:"remark"`
}

// StockInRequest 入库请求
type StockInRequest struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	UnitCost    float64 `json:"unit_cost"`
	ReferenceNo string  `json:"reference_no"`
	Notes       string  `json:"notes"`
}

// StockOutRequest 出库请求
type StockOutRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	ReferenceNo string `json:"reference_no"`
	Notes       string `json:"notes"`
}

// StockAdjustRequest 库存调整请求
type StockAdjustRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	NewQuantity int    `json:"new_quantity" binding:"required,min=0"`
	Notes       string `json:"notes"`
}

// ---------- 通用响应结构体 ----------

// Response 统一 API 响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedData 分页数据
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginationQuery 分页查询参数
type PaginationQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status"`
}

// GetOffset 计算偏移量
func (p *PaginationQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return (p.Page - 1) * p.PageSize
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalProducts   int64   `json:"total_products"`
	TotalCategories int64   `json:"total_categories"`
	TotalSuppliers  int64   `json:"total_suppliers"`
	TotalStockValue float64 `json:"total_stock_value"`
	LowStockCount   int64   `json:"low_stock_count"`
	TodayStockIn    int     `json:"today_stock_in"`
	TodayStockOut   int     `json:"today_stock_out"`
	TodayRecords    int64   `json:"today_records"`
}

// ChartData 图表数据
type ChartData struct {
	StockMovement []DailyMovement `json:"stock_movement"`
	TopProducts   []ProductRank   `json:"top_products"`
	CategoryStats []CategoryStat  `json:"category_stats"`
}

// DailyMovement 每日出入库数据
type DailyMovement struct {
	Date     string `json:"date"`
	StockIn  int    `json:"stock_in"`
	StockOut int    `json:"stock_out"`
}

// ProductRank 商品排名
type ProductRank struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// CategoryStat 分类统计
type CategoryStat struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
