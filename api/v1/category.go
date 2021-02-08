package v1

import (
	"Blog/model"
	"Blog/util/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加类别
func AddCategory(c *gin.Context) {
	var cate model.Category
	err := c.ShouldBindJSON(&cate)
	// 这里如果缺少required的字段会出err，但仍然会添加到数据库中，缺省默认值
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	code = model.ExistsCategory(&cate)
	if code == errmsg.SUCCESS {
		model.InsertCategory(&cate)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   &cate,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 查询单个类别

// 查询类别列表
func GetCategorys(c *gin.Context) {
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
	users, code := model.GetCategorysList(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   users,
		"msg":    errmsg.GetErrMsg(code),
	})

}

// 编辑类别
func UpdateCategory(c *gin.Context) {
	var user model.Category
	id, err := strconv.Atoi(c.Param("id"))
	err = c.ShouldBindJSON(&user)
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	user.ID = uint(id)
	code := model.ExistsCategory(&user)
	if code != errmsg.ERROR_CATEGORY_EXIST {
		code = model.EditCategory(&user)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 删除类别
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	code := model.DeleteCategory(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}
