package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int           // 桶的容量
	tokens   int           // 当前桶中的令牌数
	rate     int           // 令牌放入速度，单位：令牌/秒
	lastTime time.Time     // 上一次放入令牌的时间
	interval time.Duration // 令牌放入间隔
	mutex    sync.Mutex    // 互斥锁，保证并发安全
}

// NewTokenBucket 创建一个新的令牌桶
func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		rate:     rate,
		lastTime: time.Now(),
		interval: time.Second / time.Duration(rate),
	}
}

// Take 尝试从令牌桶中获取一个令牌，如果获取成功返回 true，否则返回 false
func (tb *TokenBucket) Take() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()

	// 计算应该放入的令牌数
	delta := int(now.Sub(tb.lastTime) / tb.interval)
	tb.tokens = tb.tokens + delta
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	// 更新上一次放入令牌的时间
	tb.lastTime = tb.lastTime.Add(time.Duration(delta) * tb.interval)

	// 尝试获取令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// LimiterMiddleware 返回一个 Gin 中间件函数，用于限制同一个 IP 地址在一秒内的请求次数
func LimiterMiddleware(tb *TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)

		// 尝试获取令牌，如果成功则继续处理请求，否则返回 429 Too Many Requests
		if tb.Take() {
			c.Next()
		} else {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status_code": -1,
				"status_msg":  ip + "访问过于频繁",
			})
		}
	}
}
