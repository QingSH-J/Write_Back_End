package Xuans_Gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func XuansGin() {
	r := gin.Default()

	// 解决跨域问题
	r.Use(corsMiddleware())

	// API 路由
	api := r.Group("/api")
	{
		// 测试端点
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// 用户相关接口
		api.GET("/users", getUsers)
		api.GET("/users/:id", getUserByID)
		api.POST("/users", createUser)

		// TODO: 添加更多 API 路由
	}

	// 静态文件服务（可选）
	r.StaticFS("/static", http.Dir("./static"))

	// 启动服务器
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

// CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 模拟数据
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: "1", Name: "卿诗雨", Age: 28},
	{ID: "2", Name: "唐濯枝", Age: 28},
}

// API 处理函数
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
}

func createUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 简单生成ID（实际应用中可能需要更复杂的逻辑）
	newUser.ID = fmt.Sprintf("%d", len(users)+1)

	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}
