// Package repository 提供数据访问层，封装所有数据库操作
package repository

import (
	"fmt"
	"time"

	"go-cargo/internal/models"

	"gorm.io/gorm"
)

// Repository 数据仓库，封装所有数据库操作
type Repository struct {
	db *gorm.DB
}

// New 创建 Repository 实例
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ==================== 用户 ====================

// CreateUser 创建用户
func (r *Repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByUsername 根据用户名查找用户
func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID查找用户
func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *Repository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// ==================== 分类 ====================

// ListCategories 获取分类列表
func (r *Repository) ListCategories(query *models.PaginationQuery) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	db := r.db.Model(&models.Category{})

	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	db.Count(&total)
	err := db.Order("sort_order ASC, id ASC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&categories).Error

	// 填充商品数量
	for i := range categories {
		r.db.Model(&models.Product{}).Where("category_id = ?", categories[i].ID).Count(&categories[i].ProductCount)
	}

	return categories, total, err
}

// GetAllCategories 获取所有启用的分类 (用于下拉选择)
func (r *Repository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("status = 1").Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

// GetCategoryByID 根据ID查找分类
func (r *Repository) GetCategoryByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := r.db.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// CreateCategory 创建分类
func (r *Repository) CreateCategory(cat *models.Category) error {
	return r.db.Create(cat).Error
}

// UpdateCategory 更新分类
func (r *Repository) UpdateCategory(cat *models.Category) error {
	return r.db.Save(cat).Error
}

// DeleteCategory 软删除分类
func (r *Repository) DeleteCategory(id uint) error {
	// 检查是否有关联商品
	var count int64
	r.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("该分类下有 %d 个商品，无法删除", count)
	}
	return r.db.Delete(&models.Category{}, id).Error
}

// ==================== 供应商 ====================

// ListSuppliers 获取供应商列表
func (r *Repository) ListSuppliers(query *models.PaginationQuery) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	db := r.db.Model(&models.Supplier{})

	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR contact_person LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	db.Count(&total)
	err := db.Order("id DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&suppliers).Error

	// 填充商品数量
	for i := range suppliers {
		r.db.Model(&models.Product{}).Where("supplier_id = ?", suppliers[i].ID).Count(&suppliers[i].ProductCount)
	}

	return suppliers, total, err
}

// GetAllSuppliers 获取所有启用的供应商 (用于下拉选择)
func (r *Repository) GetAllSuppliers() ([]models.Supplier, error) {
	var suppliers []models.Supplier
	err := r.db.Where("status = 1").Order("name ASC").Find(&suppliers).Error
	return suppliers, err
}

// GetSupplierByID 根据ID查找供应商
func (r *Repository) GetSupplierByID(id uint) (*models.Supplier, error) {
	var sup models.Supplier
	err := r.db.First(&sup, id).Error
	if err != nil {
		return nil, err
	}
	return &sup, nil
}

// CreateSupplier 创建供应商
func (r *Repository) CreateSupplier(sup *models.Supplier) error {
	return r.db.Create(sup).Error
}

// UpdateSupplier 更新供应商
func (r *Repository) UpdateSupplier(sup *models.Supplier) error {
	return r.db.Save(sup).Error
}

// DeleteSupplier 软删除供应商
func (r *Repository) DeleteSupplier(id uint) error {
	var count int64
	r.db.Model(&models.Product{}).Where("supplier_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("该供应商下有 %d 个商品，无法删除", count)
	}
	return r.db.Delete(&models.Supplier{}, id).Error
}

// ==================== 商品 ====================

// ListProducts 获取商品列表 (支持分页、搜索、筛选)
func (r *Repository) ListProducts(query *models.PaginationQuery, categoryID, supplierID *uint) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := r.db.Model(&models.Product{})

	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR sku LIKE ? OR barcode LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if categoryID != nil && *categoryID > 0 {
		db = db.Where("category_id = ?", *categoryID)
	}
	if supplierID != nil && *supplierID > 0 {
		db = db.Where("supplier_id = ?", *supplierID)
	}

	db.Count(&total)
	err := db.Preload("Category").Preload("Supplier").
		Order("id DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&products).Error

	return products, total, err
}

// GetProductByID 根据ID查找商品 (含关联)
func (r *Repository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Supplier").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductBySKU 根据SKU查找商品
func (r *Repository) GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// CreateProduct 创建商品
func (r *Repository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

// UpdateProduct 更新商品
func (r *Repository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

// DeleteProduct 软删除商品
func (r *Repository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// GetLowStockProducts 获取低库存商品
func (r *Repository) GetLowStockProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("current_stock <= min_stock AND min_stock > 0 AND status = 1").
		Preload("Category").
		Order("current_stock ASC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// ==================== 库存记录 ====================

// CreateInventoryRecord 创建库存记录
func (r *Repository) CreateInventoryRecord(record *models.InventoryRecord) error {
	return r.db.Create(record).Error
}

// UpdateProductStock 更新商品库存数量
func (r *Repository) UpdateProductStock(productID uint, newStock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).
		Update("current_stock", newStock).Error
}

// StockOperation 库存操作 (事务)
func (r *Repository) StockOperation(record *models.InventoryRecord, newStock int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新商品库存
		if err := tx.Model(&models.Product{}).Where("id = ?", record.ProductID).
			Update("current_stock", newStock).Error; err != nil {
			return err
		}
		// 创建操作记录
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return nil
	})
}

// ListInventoryRecords 查询库存操作记录
func (r *Repository) ListInventoryRecords(query *models.PaginationQuery, productID *uint, recordType string, startDate, endDate string) ([]models.InventoryRecord, int64, error) {
	var records []models.InventoryRecord
	var total int64

	db := r.db.Model(&models.InventoryRecord{})

	if productID != nil && *productID > 0 {
		db = db.Where("product_id = ?", *productID)
	}
	if recordType != "" {
		db = db.Where("type = ?", recordType)
	}
	if startDate != "" {
		db = db.Where("created_at >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		db = db.Where("created_at <= ?", endDate+" 23:59:59")
	}
	if query.Keyword != "" {
		db = db.Where("reference_no LIKE ? OR notes LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	db.Count(&total)
	err := db.Preload("Product").
		Order("created_at DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&records).Error

	return records, total, err
}

// ==================== 仪表盘统计 ====================

// GetDashboardStats 获取仪表盘统计
func (r *Repository) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	// 总商品数
	r.db.Model(&models.Product{}).Where("status = 1").Count(&stats.TotalProducts)

	// 总分类数
	r.db.Model(&models.Category{}).Where("status = 1").Count(&stats.TotalCategories)

	// 总供应商数
	r.db.Model(&models.Supplier{}).Where("status = 1").Count(&stats.TotalSuppliers)

	// 库存总价值
	r.db.Model(&models.Product{}).Where("status = 1").
		Select("COALESCE(SUM(current_stock * cost_price), 0)").
		Scan(&stats.TotalStockValue)

	// 低库存预警数
	r.db.Model(&models.Product{}).
		Where("current_stock <= min_stock AND min_stock > 0 AND status = 1").
		Count(&stats.LowStockCount)

	// 今日统计
	today := time.Now().Format("2006-01-02")
	todayStart := today + " 00:00:00"
	todayEnd := today + " 23:59:59"

	// 今日入库数量
	r.db.Model(&models.InventoryRecord{}).
		Where("type = ? AND created_at BETWEEN ? AND ?", models.StockIn, todayStart, todayEnd).
		Select("COALESCE(SUM(quantity), 0)").Scan(&stats.TodayStockIn)

	// 今日出库数量
	r.db.Model(&models.InventoryRecord{}).
		Where("type = ? AND created_at BETWEEN ? AND ?", models.StockOut, todayStart, todayEnd).
		Select("COALESCE(SUM(quantity), 0)").Scan(&stats.TodayStockOut)

	// 今日操作记录数
	r.db.Model(&models.InventoryRecord{}).
		Where("created_at BETWEEN ? AND ?", todayStart, todayEnd).
		Count(&stats.TodayRecords)

	return stats, nil
}

// GetChartData 获取图表数据
func (r *Repository) GetChartData() (*models.ChartData, error) {
	data := &models.ChartData{}

	// 最近30天出入库趋势
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		movement := models.DailyMovement{Date: dateStr}

		r.db.Model(&models.InventoryRecord{}).
			Where("type = ? AND DATE(created_at) = ?", models.StockIn, dateStr).
			Select("COALESCE(SUM(quantity), 0)").Scan(&movement.StockIn)

		r.db.Model(&models.InventoryRecord{}).
			Where("type = ? AND DATE(created_at) = ?", models.StockOut, dateStr).
			Select("COALESCE(SUM(quantity), 0)").Scan(&movement.StockOut)

		data.StockMovement = append(data.StockMovement, movement)
	}

	// 库存价值 TOP 10 商品
	var products []models.Product
	r.db.Where("status = 1").
		Order("(current_stock * cost_price) DESC").
		Limit(10).
		Find(&products)
	for _, p := range products {
		data.TopProducts = append(data.TopProducts, models.ProductRank{
			Name:  p.Name,
			Value: float64(p.CurrentStock) * p.CostPrice,
		})
	}

	// 各分类商品数量
	var categories []models.Category
	r.db.Where("status = 1").Order("sort_order ASC").Find(&categories)
	for _, cat := range categories {
		var count int64
		r.db.Model(&models.Product{}).Where("category_id = ? AND status = 1", cat.ID).Count(&count)
		data.CategoryStats = append(data.CategoryStats, models.CategoryStat{
			Name:  cat.Name,
			Count: count,
		})
	}

	return data, nil
}
