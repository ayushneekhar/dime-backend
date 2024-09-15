package main

import (
	"fmt"
	"time"

	"github.com/ayushneekhar/dime-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("expenses.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	DB.AutoMigrate(&models.Transaction{}, &models.Category{})
}

func scheduleRecurringTransactions() {
	c := cron.New()
	c.AddFunc(("@daily"), func() {
		var recurringTransactions []models.Transaction
		DB.Where("recurring = ?", true).Find(&recurringTransactions)
		for _, transaction := range recurringTransactions {
			transaction.ID = 0
			transaction.Timestamp = time.Now()
			if err := DB.Create(&transaction).Error; err != nil {
				fmt.Println("Error creating recurring transaction", err)
			}
		}
	})
	c.Start()
}

func main() {
	initDatabase()
	scheduleRecurringTransactions()
	r := gin.Default()

	r.POST("/category", CreateCategory)
	r.GET("/category", GetCategories)
	r.POST("/transaction", CreateTransaction)
	r.GET("/transaction", GetTransactions)

	r.Run()
}

func CreateCategory(c *gin.Context) {
	var category models.Category // Should be models.Category, not []models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Create(&category).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := DB.Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, categories)
}

func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindBodyWithJSON(&transaction); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, transaction)
}

func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	if err := DB.Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, transactions)
}
