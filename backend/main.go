package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-jwt-webdemo/claim"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Secret    string `yaml:"secret"`
	RedisAddr string `yaml:"redisAddr"`
}

var client *redis.Client
var secret []byte

func init() {
	/*	LoadConfig	*/
	file, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*	InitRedis	*/
	client = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr, // 你的Redis服务器地址和端口
	})

	_, err = client.Ping().Result()
	if err != nil {
		fmt.Println("连接Redis失败:", err)
		return
	}
	/*	InitSecret	*/
	secret = []byte(config.Secret)
}

func main() {
	r := gin.Default()

	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)

	api := r.Group("/api")
	api.Use(jwtAuthMiddleware())
	api.POST("/order", orderHandler)

	r.Run(":9099")
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerHandler(c *gin.Context) {
	var req LoginReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	// 检查用户是否已经存在
	exists, err := client.Exists(req.Username).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	if exists == 1 {
		c.JSON(http.StatusConflict, "user already exists")
		return
	}

	// 在Redis中存储用户信息
	err = client.Set(req.Username, req.Password, time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	c.JSON(http.StatusOK, "user registered successfully")
}

func loginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	// 从Redis中获取用户信息
	password, err := client.Get(req.Username).Result()
	if err != nil {
		c.JSON(http.StatusForbidden, "invalid username or password")
		return
	}

	if password == req.Password {
		token, err := claim.GenToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
		return
	}

	c.JSON(http.StatusForbidden, "invalid username or password")
}

type OrderReq struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
}

func orderHandler(c *gin.Context) {
	var req OrderReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	greet := fmt.Sprintf("Hi %v, I will give you %v %v", userId, req.Count, req.Product)
	c.JSON(http.StatusOK, greet)
}

func jwtAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusForbidden, "empty token")
			c.Abort()
			return
		}
		claim, err := claim.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, "Invalid token")
			c.Abort()
			return
		}
		c.Set("userId", claim.UserName)
		c.Next()
	}
}
