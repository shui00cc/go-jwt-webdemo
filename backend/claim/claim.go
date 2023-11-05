package claim

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const TokenExpireDuration = 5 * time.Minute

type Config struct {
	Secret        string `yaml:"secret"`
	RedisAddr     string `yaml:"redisAddr"`
	Authorization string `yaml:"authorization"`
}

/*	全局变量	*/
var Client *redis.Client
var Authorization string // 用于AI绘画的api token: https://open.creator.nolibox.com/
/*	----	*/
var secret []byte // 用于生成登录token的密钥

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
	Client = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr, // 你的Redis服务器地址和端口
	})

	_, err = Client.Ping().Result()
	if err != nil {
		fmt.Println("连接Redis失败:", err)
		return
	}
	/*	InitOthers	*/
	secret = []byte(config.Secret)
	Authorization = config.Authorization
}

type AuthClaim struct {
	UserName string `json:"userName"`
	jwt.StandardClaims
	Authorization string
}

func GenToken(userName string) (string, error) {
	c := AuthClaim{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "CC",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*AuthClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaim{}, func(tk *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claim, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("Invalid token ")
}
