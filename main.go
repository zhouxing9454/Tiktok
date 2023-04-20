package main

import (
	"TikTok_Project/repository"
	"TikTok_Project/router"
	"TikTok_Project/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := repository.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := repository.Close()
		if err != nil {
			log.Println("can't close current db！")
		}
	}()
	repository.ModelAutoMigrate()

	if err := repository.InitRedisClient(); err != nil {
		panic(err)
	}

	if err := utils.SensitiveWordInit(); err != nil {
		log.Printf("敏感词初始化失败")
		panic(err)
	}

	r := router.InitRouter()
	srv := &http.Server{
		Addr:    "IP:PORT",
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
