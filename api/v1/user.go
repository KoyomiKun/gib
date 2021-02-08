package v1

import (
	"Blog/middleware"
	"Blog/model"
	"Blog/util"
	"Blog/util/errmsg"
	"Blog/util/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var code errmsg.Code
var logger = util.GetLogger()

// 普通用户注册
func Register(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	// 这里如果缺少required的字段会出err，但仍然会添加到数据库中，缺省默认值
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	msg, code := validator.Validate(&u)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    msg,
		})
		return
	}
	_, code = model.ExistsUser(u.Username)
	if code == errmsg.ERROR_USER_NOT_EXIST {
		u.Role = 2
		model.InsertUser(&u)
		code = errmsg.SUCCESS
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"data":   &u,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	} else {
		code = errmsg.ERROR_USERNAME_USED
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
}

// 添加管理员用户
func AddAdmin(c *gin.Context) {
	var u model.User
	role := c.GetInt("role")
	// 只有管理员能创建管理员
	if role != 1 {
		code = errmsg.ERROR_NOT_ADMIN
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	err := c.ShouldBindJSON(&u)
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	_, code = model.ExistsUser(u.Username)
	if code == errmsg.ERROR_USER_NOT_EXIST {
		u.Role = 1
		model.InsertUser(&u)
		code = errmsg.SUCCESS
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"data":   &u,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	} else {
		code = errmsg.ERROR_USERNAME_USED
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
}

// 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Query("pagesize"))
	pageNum, err := strconv.Atoi(c.Query("pagenum"))
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}

	if pageSize == 0 {
		// 如果limit是-1就是不要限制
		pageSize = -1
	}
	if pageNum == 0 {
		// 如果offset为-1就是不要偏移
		pageNum = -1
	}
	users, code := model.GetUsersList(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   users,
		"msg":    errmsg.GetErrMsg(code),
	})

}

// 编辑个人信息
func UpdateUser(c *gin.Context) {
	var user model.User
	//id, err := strconv.Atoi(c.Param("id"))
	id := c.GetInt("id")
	err := c.ShouldBindJSON(&user)
	if err != nil || id == 0 {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	user.ID = uint(id)
	_, code := model.ExistsUser(user.Username)
	if code != errmsg.SUCCESS {
		code = model.EditUser(user.ID, &user)
	} else {
		code = errmsg.ERROR_USERNAME_USED
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	cRole := c.GetInt("role")
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	// 非管理员无法删除用户
	if cRole != 1 {
		code = errmsg.ERROR_NOT_ADMIN
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	code := model.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	if code = model.Login(user.Username, user.Password); code == errmsg.SUCCESS {
		id, _ := model.ExistsUser(user.Username)
		u, _ := model.SelectUser(id)
		token, code := middleware.GetJwt().CreateToken(u.Username, u.ID, u.Role)
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
			"token":  token,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

}

// 上传头像
func UploadAvatar(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	username := c.GetString("username")
	if id, code := model.ExistsUser(username); code == errmsg.SUCCESS {
		fileSize := fileHeader.Size
		url, _ := util.UploadFile(file, fileSize)
		code = model.EditUser(uint(id), &model.User{Avatar: url})
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
			"url":    url,
		})
		return
	} else {
		code = errmsg.ERROR_USER_NOT_EXIST
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

}
