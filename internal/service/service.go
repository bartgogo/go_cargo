// Package service 提供业务逻辑层
package service

import (
	"errors"
	"fmt"
	"time"

	"go-cargo/internal/config"
	"go-cargo/internal/models"
	"go-cargo/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service 业务逻辑层
type Service struct {
	repo *repository.Repository
	cfg  *config.Config
}

// New 创建 Service 实例
func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

// ==================== 认证 ====================

// Login 用户登录, 返回 JWT token
func (s *Service) Login(req *models.LoginRequest) (string, *models.User, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("用户名或密码错误")
		}
		return "", nil, fmt.Errorf("登录失败: %w", err)
	}

	if user.Status != 1 {
		return "", nil, fmt.Errorf("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, fmt.Errorf("用户名或密码错误")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return token, user, nil
}

// Register 用户注册
func (s *Service) Register(req *models.RegisterRequest) (*models.User, error) {
	// 检查用户名是否已存在
	if _, err := s.repo.GetUserByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPwd),
		Email:    req.Email,
		RealName: req.RealName,
		Role:     "operator",
		Status:   1,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// GetProfile 获取用户信息
func (s *Service) GetProfile(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// UpdateProfile 更新个人信息
func (s *Service) UpdateProfile(userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	user.Email = req.Email
	user.RealName = req.RealName
	user.Phone = req.Phone
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("更新失败: %w", err)
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("原密码错误")
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.Password = string(hashedPwd)
	return s.repo.UpdateUser(user)
}

// generateToken 生成 JWT Token
func (s *Service) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(s.cfg.JWTExpireHours) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// ==================== 分类 ====================

// ListCategories 获取分类列表
func (s *Service) ListCategories(query *models.PaginationQuery) ([]models.Category, int64, error) {
	return s.repo.ListCategories(query)
}

// GetAllCategories 获取所有启用分类
func (s *Service) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}

// CreateCategory 创建分类
func (s *Service) CreateCategory(req *models.CategoryRequest) (*models.Category, error) {
	cat := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      1,
	}
	if req.Status != 0 {
		cat.Status = req.Status
	}
	if err := s.repo.CreateCategory(cat); err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}
	return cat, nil
}

// UpdateCategory 更新分类
func (s *Service) UpdateCategory(id uint, req *models.CategoryRequest) (*models.Category, error) {
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("分类不存在")
	}
	cat.Name = req.Name
	cat.Description = req.Description
	cat.SortOrder = req.SortOrder
	if req.Status != 0 {
		cat.Status = req.Status
	}
	if err := s.repo.UpdateCategory(cat); err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}
	return cat, nil
}

// DeleteCategory 删除分类
func (s *Service) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}

// ==================== 供应商 ====================

// ListSuppliers 获取供应商列表
func (s *Service) ListSuppliers(query *models.PaginationQuery) ([]models.Supplier, int64, error) {
	return s.repo.ListSuppliers(query)
}

// GetAllSuppliers 获取所有启用供应商
func (s *Service) GetAllSuppliers() ([]models.Supplier, error) {
	return s.repo.GetAllSuppliers()
}

// CreateSupplier 创建供应商
func (s *Service) CreateSupplier(req *models.SupplierRequest) (*models.Supplier, error) {
	sup := &models.Supplier{
		Code:          req.Code,
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Remark:        req.Remark,
		Status:        1,
	}
	if req.Status != 0 {
		sup.Status = req.Status
	}
	if err := s.repo.CreateSupplier(sup); err != nil {
		return nil, fmt.Errorf("创建供应商失败: %w", err)
	}
	return sup, nil
}

// UpdateSupplier 更新供应商
func (s *Service) UpdateSupplier(id uint, req *models.SupplierRequest) (*models.Supplier, error) {
	sup, err := s.repo.GetSupplierByID(id)
	if err != nil {
		return nil, fmt.Errorf("供应商不存在")
	}
	sup.Code = req.Code
	sup.Name = req.Name
	sup.ContactPerson = req.ContactPerson
	sup.Phone = req.Phone
	sup.Email = req.Email
	sup.Address = req.Address
	sup.Remark = req.Remark
	if req.Status != 0 {
		sup.Status = req.Status
	}
	if err := s.repo.UpdateSupplier(sup); err != nil {
		return nil, fmt.Errorf("更新供应商失败: %w", err)
	}
	return sup, nil
}

// DeleteSupplier 删除供应商
func (s *Service) DeleteSupplier(id uint) error {
	return s.repo.DeleteSupplier(id)
}

// ==================== 商品 ====================

// ListProducts 获取商品列表
func (s *Service) ListProducts(query *models.PaginationQuery, categoryID, supplierID *uint) ([]models.Product, int64, error) {
	return s.repo.ListProducts(query, categoryID, supplierID)
}

// GetProduct 获取商品详情
func (s *Service) GetProduct(id uint) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

// CreateProduct 创建商品
func (s *Service) CreateProduct(req *models.ProductRequest) (*models.Product, error) {
	// 检查 SKU 是否重复
	if _, err := s.repo.GetProductBySKU(req.SKU); err == nil {
		return nil, fmt.Errorf("SKU '%s' 已存在", req.SKU)
	}

	product := &models.Product{
		SKU:          req.SKU,
		Name:         req.Name,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		SupplierID:   req.SupplierID,
		Unit:         req.Unit,
		CostPrice:    req.CostPrice,
		SellingPrice: req.SellingPrice,
		MinStock:     req.MinStock,
		MaxStock:     req.MaxStock,
		Barcode:      req.Barcode,
		Location:     req.Location,
		ImageURL:     req.ImageURL,
		Status:       1,
	}
	if req.Unit == "" {
		product.Unit = "个"
	}
	if req.Status != 0 {
		product.Status = req.Status
	}

	if err := s.repo.CreateProduct(product); err != nil {
		return nil, fmt.Errorf("创建商品失败: %w", err)
	}
	return product, nil
}

// UpdateProduct 更新商品
func (s *Service) UpdateProduct(id uint, req *models.ProductRequest) (*models.Product, error) {
	product, err := s.repo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("商品不存在")
	}

	// 如果 SKU 有变更，检查新 SKU 是否重复
	if req.SKU != product.SKU {
		if existing, err := s.repo.GetProductBySKU(req.SKU); err == nil && existing.ID != id {
			return nil, fmt.Errorf("SKU '%s' 已存在", req.SKU)
		}
	}

	product.SKU = req.SKU
	product.Name = req.Name
	product.Description = req.Description
	product.CategoryID = req.CategoryID
	product.SupplierID = req.SupplierID
	product.Unit = req.Unit
	product.CostPrice = req.CostPrice
	product.SellingPrice = req.SellingPrice
	product.MinStock = req.MinStock
	product.MaxStock = req.MaxStock
	product.Barcode = req.Barcode
	product.Location = req.Location
	product.ImageURL = req.ImageURL
	if req.Status != 0 {
		product.Status = req.Status
	}

	if err := s.repo.UpdateProduct(product); err != nil {
		return nil, fmt.Errorf("更新商品失败: %w", err)
	}
	return product, nil
}

// DeleteProduct 删除商品
func (s *Service) DeleteProduct(id uint) error {
	return s.repo.DeleteProduct(id)
}

// ==================== 库存操作 ====================

// StockIn 入库操作
func (s *Service) StockIn(req *models.StockInRequest, operatorID uint, operatorName string) error {
	product, err := s.repo.GetProductByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	beforeQty := product.CurrentStock
	afterQty := beforeQty + req.Quantity
	totalCost := float64(req.Quantity) * req.UnitCost

	record := &models.InventoryRecord{
		ProductID:    req.ProductID,
		Type:         models.StockIn,
		Quantity:     req.Quantity,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		UnitCost:     req.UnitCost,
		TotalCost:    totalCost,
		ReferenceNo:  req.ReferenceNo,
		Notes:        req.Notes,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	}

	return s.repo.StockOperation(record, afterQty)
}

// StockOut 出库操作
func (s *Service) StockOut(req *models.StockOutRequest, operatorID uint, operatorName string) error {
	product, err := s.repo.GetProductByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	if product.CurrentStock < req.Quantity {
		return fmt.Errorf("库存不足，当前库存: %d，请求出库: %d", product.CurrentStock, req.Quantity)
	}

	beforeQty := product.CurrentStock
	afterQty := beforeQty - req.Quantity

	record := &models.InventoryRecord{
		ProductID:    req.ProductID,
		Type:         models.StockOut,
		Quantity:     req.Quantity,
		BeforeQty:    beforeQty,
		AfterQty:     afterQty,
		ReferenceNo:  req.ReferenceNo,
		Notes:        req.Notes,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	}

	return s.repo.StockOperation(record, afterQty)
}

// StockAdjust 库存调整
func (s *Service) StockAdjust(req *models.StockAdjustRequest, operatorID uint, operatorName string) error {
	product, err := s.repo.GetProductByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("商品不存在")
	}

	beforeQty := product.CurrentStock
	diff := req.NewQuantity - beforeQty
	quantity := diff
	if quantity < 0 {
		quantity = -quantity
	}

	record := &models.InventoryRecord{
		ProductID:    req.ProductID,
		Type:         models.StockAdjust,
		Quantity:     quantity,
		BeforeQty:    beforeQty,
		AfterQty:     req.NewQuantity,
		Notes:        req.Notes,
		OperatorID:   operatorID,
		OperatorName: operatorName,
	}

	return s.repo.StockOperation(record, req.NewQuantity)
}

// ListInventoryRecords 查询库存记录
func (s *Service) ListInventoryRecords(query *models.PaginationQuery, productID *uint, recordType, startDate, endDate string) ([]models.InventoryRecord, int64, error) {
	return s.repo.ListInventoryRecords(query, productID, recordType, startDate, endDate)
}

// ==================== 仪表盘 ====================

// GetDashboardStats 获取仪表盘统计
func (s *Service) GetDashboardStats() (*models.DashboardStats, error) {
	return s.repo.GetDashboardStats()
}

// GetChartData 获取图表数据
func (s *Service) GetChartData() (*models.ChartData, error) {
	return s.repo.GetChartData()
}

// GetLowStockProducts 获取低库存商品
func (s *Service) GetLowStockProducts() ([]models.Product, error) {
	return s.repo.GetLowStockProducts(20)
}
