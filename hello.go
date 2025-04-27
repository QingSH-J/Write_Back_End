package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinxinyu/go_backend/Xuans_Gin"
)

func main() {
	fmt.Println("Hello, World!")

	// 在主goroutine中运行Gin服务器，而不是在新的goroutine中
	go Xuans_Gin.XuansGin()

	// 创建一个通道来接收信号
	quit := make(chan os.Signal, 1)

	// 监听SIGINT和SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到收到信号
	<-quit

	fmt.Println("\n服务器正在关闭...")
}
