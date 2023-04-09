package middleware

import (
	"TikTok_Project/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// RateMiddleware 中间件
func RateMiddleware(c *gin.Context) {
	//以Pipeline的方式操作事务
	pipe := repository.RDB.TxPipeline()

	// 5 秒刷新key为IP(c.ClientIP())的r值为0
	err := pipe.SetNX(repository.CTX, c.ClientIP(), 0, 10*time.Second).Err()
	if err != nil {
		log.Printf("redis刷新错误" + err.Error())
	}
	// 每次访问，这个IP的对应的值加一
	pipe.Incr(repository.CTX, c.ClientIP())
	// 提交事务
	_, _ = pipe.Exec(repository.CTX)

	// 获取IP访问的次数
	var val int
	val, err = repository.RDB.Get(repository.CTX, c.ClientIP()).Int()
	if err != nil {
		log.Printf("redis刷新错误" + err.Error())
	}
	// 如果10秒内大于50次
	if val > 50 {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  "访问过于频繁",
		})
	} else {
		c.Next()
	}
}
