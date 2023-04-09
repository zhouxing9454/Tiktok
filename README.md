# Tiktok
自己独立修改的抖声项目，基于青训营团队抖声项目。原项目是我们青训营团队的大项目作业，这个是原项目地址：https://github.com/Lionel24-xxy/douyin-project





### 修改1

将sha-1改为sha-256算法，加盐值，盐值和password一起存入数据库。

```Go
func HashPassword(password string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func GenerateSalt(length int) string {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(salt)
}
```





### 修改2

jwt的库修改（原来的库已经不再维护了）

```go
import "github.com/golang-jwt/jwt"
```





### 修改3

删除normal.go，将对应的中间件放入AuthMiddlerware里面
