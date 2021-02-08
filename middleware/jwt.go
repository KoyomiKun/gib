package middleware

import (
	"Blog/util"
	"Blog/util/errmsg"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var code errmsg.Code

var once sync.Once
var myJwt Jwt

type Jwt struct {
	SignKey []byte
}
type Claims struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Role     int    `json:"role,omitempty"`
	jwt.StandardClaims
}

func GetJwt() *Jwt {
	once.Do(func() {
		myJwt.SignKey = []byte(util.JwtKey)
	})
	return &myJwt
}

// 生成token
func (j *Jwt) CreateToken(username string, id uint, role int) (string, errmsg.Code) {
	// 有效时间
	expireTime := time.Now().Add(10 * time.Hour)
	c := Claims{
		id,
		username,
		role,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "komikunblog",
		},
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := claim.SignedString(j.SignKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 验证token
func (j *Jwt) VarifyToken(tokenString string) (*Claims, errmsg.Code) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if key, ok := token.Claims.(*Claims); ok && token.Valid {
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR
	}
}

// jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 规范写法
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_NOT_EXIST
			// abort不会停止现在的函数，只会停止挂起函数
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		varifyToken := strings.SplitN(tokenHeader, " ", 2)
		if len(varifyToken) != 2 && varifyToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_ILLIGAL
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		j := GetJwt()
		key, code2 := j.VarifyToken(varifyToken[1])
		if code2 == errmsg.ERROR {
			code = errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > key.ExpiresAt {
			code = errmsg.ERROR_TOKEN_OVERTIME
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		//code = errmsg.SUCCESS
		//c.JSON(http.StatusOK, gin.H{
		//"status": code,
		//"msg":    errmsg.GetErrMsg(code),
		//})
		c.Set("username", key.Username)
		c.Set("id", key.ID)
		c.Set("role", key.Role)
		c.Next()
	}
}
