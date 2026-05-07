package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- MODELS (Cấu trúc dữ liệu) ---

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}

type Order struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"` // Shopee, TikTok, Facebook, Zalo...
	Customer  string    `json:"customer"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountingSummary struct {
	Revenue int `json:"revenue"`
	Expense int `json:"expense"`
	Profit  int `json:"profit"`
}

func main() {
	// Khởi tạo Gin với chế độ Release để tối ưu hiệu năng (nếu muốn)
	r := gin.Default()

	// MIDDLEWARE: Xử lý CORS để React (Port 5173) gọi được API (Port 8080)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- DATA MẪU (Mock Data) ---
	
	inventory := []Product{
		{ID: 1, Name: "Thức ăn mèo Whiskas", Stock: 120, Price: 150000},
		{ID: 2, Name: "Cát vệ sinh Đậu Nành", Stock: 45, Price: 90000},
		{ID: 3, Name: "Pate King's Pet", Stock: 200, Price: 35000},
	}

	allOrders := []Order{
		{ID: "SHP-001", Source: "shopee", Customer: "Minh Nguyễn", Amount: 300000, CreatedAt: time.Now()},
		{ID: "TT-992", Source: "tiktok", Customer: "Lê An", Amount: 150000, CreatedAt: time.Now()},
		{ID: "FB-123", Source: "facebook", Customer: "Hoàng Oanh", Amount: 450000, CreatedAt: time.Now()},
		{ID: "YT-55", Source: "youtube", Customer: "Quốc Anh", Amount: 120000, CreatedAt: time.Now()},
		{ID: "ZLO-09", Source: "zalo", Customer: "Bảo Thy", Amount: 210000, CreatedAt: time.Now()},
	}

	// --- API ENDPOINTS ---

	api := r.Group("/api")
	{
		// 1. API Kho hàng
		api.GET("/inventory", func(c *gin.Context) {
			c.JSON(http.StatusOK, inventory)
		})

		// 2. API Đơn hàng theo từng sàn (Sử dụng Param :source)
		api.GET("/orders/:source", func(c *gin.Context) {
			source := c.Param("source")
			var filteredOrders []Order
			
			for _, order := range allOrders {
				if order.Source == source {
					filteredOrders = append(filteredOrders, order)
				}
			}
			c.JSON(http.StatusOK, filteredOrders)
		})

		// 3. API Kế toán
		api.GET("/accounting", func(c *gin.Context) {
			summary := AccountingSummary{
				Revenue: 125000000, // 125 triệu
				Expense: 85000000,  // 85 triệu
				Profit:  40000000,  // 40 triệu
			}
			c.JSON(http.StatusOK, summary)
		})

		// 4. API Lên đơn hàng mới (Ví dụ tích hợp từ tin nhắn Zalo)
		api.POST("/orders", func(c *gin.Context) {
			var newOrder Order
			if err := c.ShouldBindJSON(&newOrder); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đơn hàng không hợp lệ"})
				return
			}
			newOrder.CreatedAt = time.Now()
			// Ở đây bạn sẽ viết logic lưu vào Database (MySQL/PostgreSQL)
			allOrders = append(allOrders, newOrder)
			
			c.JSON(http.StatusCreated, gin.H{
				"status": "success",
				"message": "Đã ghi nhận đơn hàng từ " + newOrder.Source,
				"data": newOrder,
			})
		})
	}

	// Chạy server tại cổng 8080
	r.Run(":8080")
}