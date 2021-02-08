package v1

import (
	"Blog/model"
	"Blog/util/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 添加文章
func AddPost(c *gin.Context) {
	var p model.Post
	err := c.ShouldBindJSON(&p)
	// 这里如果缺少required的字段会出err，但仍然会添加到数据库中，缺省默认值
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	code = model.InsertPost(&p)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   &p,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 查询单个文章
func GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	post, code := model.GetPost(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   post,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 查询类别下文章
func GetPostByCate(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Query("pagesize"))
	pageNum, err := strconv.Atoi(c.Query("pagenum"))
	cid, err := strconv.Atoi(c.Param("cid"))
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
	posts, code := model.GetPostByCate(cid, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data": gin.H{
			"count": len(posts),
			"posts": posts,
		},
		"msg": errmsg.GetErrMsg(code),
	})
}

// 查询文章列表
func GetPosts(c *gin.Context) {
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
	posts, code := model.GetPostsList(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data": gin.H{
			"count": len(posts),
			"posts": posts,
		},
		"msg": errmsg.GetErrMsg(code),
	})

}

// 编辑文章
func UpdatePost(c *gin.Context) {
	var post model.Post
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}
	c.ShouldBindJSON(&post)
	post.ID = uint(id)
	code = model.EditPost(&post)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// 删除文章
func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		code = errmsg.PARAM_ERROR
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code) + err.Error(),
		})
		return
	}

	code := model.DeletePost(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    errmsg.GetErrMsg(code),
	})
}
