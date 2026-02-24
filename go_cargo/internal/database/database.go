// Package database 提供数据库初始化和连接管理
package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"go-cargo/internal/config"
	"go-cargo/internal/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// Init 初始化数据库连接、自动迁移、种子数据
func Init(cfg *config.Config) *gorm.DB {
	// 确保数据目录存在
	dbDir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("[DB] 创建数据目录失败: %v", err)
	}

	// 配置 GORM 日志
	logLevel := logger.Info
	if cfg.AppMode == "release" {
		logLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接 SQLite
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
	if err != nil {
		log.Fatalf("[DB] 连接数据库失败: %v", err)
	}

	// 启用 WAL 模式提升并发性能
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		log.Fatalf("[DB] 自动迁移失败: %v", err)
	}

	// 初始化种子数据
	seedData(db, cfg)

	DB = db
	log.Println("[DB] 数据库初始化完成")
	return db
}

// autoMigrate 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Supplier{},
		&models.Product{},
		&models.InventoryRecord{},
	)
}

// seedData 初始化种子数据 (仅首次运行)
func seedData(db *gorm.DB, cfg *config.Config) {
	// 创建默认管理员
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		admin := models.User{
			Username: cfg.AdminUsername,
			Password: string(hashedPwd),
			RealName: "系统管理员",
			Role:     "admin",
			Status:   1,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("[DB] 创建默认管理员失败: %v", err)
		} else {
			log.Printf("[DB] 默认管理员已创建: %s / %s", cfg.AdminUsername, cfg.AdminPassword)
		}
	}

	// 创建示例分类
	var catCount int64
	db.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		categories := []models.Category{
			{Name: "电子产品", Description: "手机、电脑、配件等电子类商品", SortOrder: 1, Status: 1},
			{Name: "办公用品", Description: "文具、纸张、打印耗材等", SortOrder: 2, Status: 1},
			{Name: "食品饮料", Description: "零食、饮品、生鲜等", SortOrder: 3, Status: 1},
			{Name: "服装鞋帽", Description: "男装、女装、童装、鞋类等", SortOrder: 4, Status: 1},
			{Name: "家居日用", Description: "家具、家纺、清洁用品等", SortOrder: 5, Status: 1},
		}
		db.Create(&categories)
		log.Println("[DB] 示例分类数据已创建")
	}

	// 创建示例供应商
	var supCount int64
	db.Model(&models.Supplier{}).Count(&supCount)
	if supCount == 0 {
		suppliers := []models.Supplier{
			{Code: "SUP001", Name: "深圳科技有限公司", ContactPerson: "张经理", Phone: "13800138001", Email: "zhang@example.com", Address: "广东省深圳市南山区科技园", Status: 1},
			{Code: "SUP002", Name: "上海贸易集团", ContactPerson: "李总", Phone: "13800138002", Email: "li@example.com", Address: "上海市浦东新区陆家嘴", Status: 1},
			{Code: "SUP003", Name: "北京文具厂", ContactPerson: "王主管", Phone: "13800138003", Email: "wang@example.com", Address: "北京市朝阳区望京", Status: 1},
		}
		db.Create(&suppliers)
		log.Println("[DB] 示例供应商数据已创建")
	}

	// 创建示例商品
	var prodCount int64
	db.Model(&models.Product{}).Count(&prodCount)
	if prodCount == 0 {
		catID1 := uint(1)
		catID2 := uint(2)
		supID1 := uint(1)
		supID2 := uint(2)
		supID3 := uint(3)
		products := []models.Product{
			{SKU: "P001", Name: "MacBook Pro 14寸", CategoryID: &catID1, SupplierID: &supID1, Unit: "台", CostPrice: 12000, SellingPrice: 15999, CurrentStock: 25, MinStock: 5, Location: "A-01-01", Status: 1},
			{SKU: "P002", Name: "iPhone 16 Pro", CategoryID: &catID1, SupplierID: &supID1, Unit: "台", CostPrice: 6000, SellingPrice: 8999, CurrentStock: 50, MinStock: 10, Location: "A-01-02", Status: 1},
			{SKU: "P003", Name: "机械键盘 K8 Pro", CategoryID: &catID1, SupplierID: &supID2, Unit: "个", CostPrice: 350, SellingPrice: 599, CurrentStock: 120, MinStock: 20, Location: "A-02-01", Status: 1},
			{SKU: "P004", Name: "A4 打印纸 (500张/包)", CategoryID: &catID2, SupplierID: &supID3, Unit: "包", CostPrice: 18, SellingPrice: 28, CurrentStock: 200, MinStock: 50, Location: "B-01-01", Status: 1},
			{SKU: "P005", Name: "中性笔 (黑色 0.5mm)", CategoryID: &catID2, SupplierID: &supID3, Unit: "支", CostPrice: 1.5, SellingPrice: 3, CurrentStock: 500, MinStock: 100, Location: "B-01-02", Status: 1},
			{SKU: "P006", Name: "无线鼠标 M720", CategoryID: &catID1, SupplierID: &supID2, Unit: "个", CostPrice: 180, SellingPrice: 299, CurrentStock: 80, MinStock: 15, Location: "A-02-02", Status: 1},
			{SKU: "P007", Name: "显示器支架", CategoryID: &catID1, SupplierID: &supID2, Unit: "个", CostPrice: 120, SellingPrice: 199, CurrentStock: 3, MinStock: 5, Location: "A-03-01", Status: 1},
			{SKU: "P008", Name: "USB-C 扩展坞", CategoryID: &catID1, SupplierID: &supID1, Unit: "个", CostPrice: 200, SellingPrice: 359, CurrentStock: 2, MinStock: 10, Location: "A-03-02", Status: 1},
		}
		db.Create(&products)
		log.Println("[DB] 示例商品数据已创建")
	}
}
