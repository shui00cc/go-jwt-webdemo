package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-jwt-webdemo/claim"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "token", "Authorization"},
	}))
	r.POST("/login", loginHandler)
	r.POST("/register", registerHandler)

	api := r.Group("/api")
	api.Use(jwtAuthMiddleware())
	api.POST("/order", orderHandler)
	api.POST("/config", configHandler)

	r.Run(":9099")
}

/*	以下是登录注册处理	*/
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// encryptPassword 对密码进行加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte("ccsecret"))
	return hex.EncodeToString(h.Sum(data))
}

func registerHandler(c *gin.Context) {
	var req LoginReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "internal error",
		})
		return
	}

	// 检查用户是否已经存在
	exists, err := claim.Client.Exists(req.Username).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	if exists == 1 {
		c.JSON(http.StatusConflict, "user already exists")
		return
	}

	// 在Redis中存储用户信息
	req.Password = encryptPassword([]byte(req.Password))
	err = claim.Client.Set(req.Username, req.Password, time.Hour).Err()
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
	password, err := claim.Client.Get(req.Username).Result()
	if err != nil {
		c.JSON(http.StatusForbidden, "invalid username or password")
		return
	}

	if password == encryptPassword([]byte(req.Password)) {
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

/*	以下是order处理	*/
type OrderReq struct {
	Text string `json:"text"`
}

type DrawRequest struct {
	Task   string                 `json:"task"`
	Params map[string]interface{} `json:"params"`
}

type DrawResponse struct {
	UID     string `json:"uid"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

func orderHandler(c *gin.Context) {
	var req OrderReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	drawReq := DrawRequest{
		Task: "txt2img.sd",
		Params: map[string]interface{}{
			"text":     req.Text,
			"w":        512,
			"h":        512,
			"is_anime": false,
		},
	}

	reqBody, err := json.Marshal(drawReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	apiReq, err := http.NewRequest("POST", "https://open.nolibox.com/prod-open-aigc/engine/push", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}

	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", claim.Authorization)

	client := &http.Client{}
	resp, err := client.Do(apiReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	}
	defer resp.Body.Close()

	// 处理响应
	var drawRes DrawResponse
	if err := json.NewDecoder(resp.Body).Decode(&drawRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "internal error"})
		return
	} else {
		c.JSON(http.StatusOK, drawRes)
	}
}

/*	以下是config处理	*/
type Config struct {
	Secret        string `yaml:"secret"`
	RedisAddr     string `yaml:"redisAddr"`
	Authorization string `yaml:"authorization"`
}

func configHandler(c *gin.Context) {
	config, err := readConfig("config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authorization": config.Authorization,
	})
}

func readConfig(filename string) (Config, error) {
	var config Config

	// 读取config.yaml文件
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// 解析yaml
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

/*	中间件	*/
func jwtAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusForbidden, "empty token")
			c.Abort()
			return
		}
		authClaim, err := claim.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, "Invalid token")
			c.Abort()
			return
		}
		c.Set("userName", authClaim.UserName)
		c.Next()
	}
}
