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
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	configInit(wg)
	wg.Wait()

	defer func() {
		err := repository.Close()
		if err != nil {
			log.Println("can't close current db！")
		}
	}()

	r := router.InitRouter()
	srv := &http.Server{
		Addr:    ":8000", //自定义
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

func configInit(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		err := repository.InitMySQL()
		if err != nil {
			panic(err)
		}
		repository.ModelAutoMigrate()

	}()
	go func() {
		defer wg.Done()
		if err := repository.InitRedisClient(); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := utils.SensitiveWordInit(); err != nil {
			log.Printf("敏感词初始化失败")
			panic(err)
		}
	}()
}
