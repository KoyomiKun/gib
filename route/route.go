package route

import (
	v1 "Blog/api/v1"
	"Blog/middleware"
	"Blog/util"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	gin.SetMode(util.AppMode)
	// 和new的区别就是加了两个中间件
	//不要他的日志，自己写
	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger())
	// 鉴权
	routeV1Auth := r.Group("api/v1")
	routeV1Auth.Use(middleware.JwtToken())
	{
		//User 模块路由接口
		routeUser := routeV1Auth.Group("user")
		{
			routeUser.POST("/admin", v1.AddAdmin)
			routeUser.PUT("/", v1.UpdateUser)
			routeUser.DELETE("/:id", v1.DeleteUser)
			routeUser.POST("/avatar", v1.UploadAvatar)
		}
		//Category 模块路由接口
		routeCategory := routeV1Auth.Group("category")
		{
			routeCategory.POST("/", v1.AddCategory)
			routeCategory.PUT("/:id", v1.UpdateCategory)
			routeCategory.DELETE("/:id", v1.DeleteCategory)
		}
		//Post 模块路由接口
		routePost := routeV1Auth.Group("post")
		{
			routePost.POST("/", v1.AddPost)
			routePost.PUT("/:id", v1.UpdatePost)
			routePost.DELETE("/:id", v1.DeletePost)
		}
	}
	// 公共
	routeV1Public := r.Group("api/v1")
	{
		//User 模块路由接口
		routeUser := routeV1Public.Group("user")
		{
			routeUser.POST("/", v1.Register)
			routeUser.GET("/", v1.GetUsers) // 之后要放到认证中去
		}
		//Category 模块路由接口
		routeCategory := routeV1Public.Group("category")
		{
			routeCategory.GET("/", v1.GetCategorys)
		}
		//Post 模块路由接口
		routePost := routeV1Public.Group("post")
		{
			routePost.GET("/posts", v1.GetPosts)
			routePost.GET("/post/:id", v1.GetPost)
			routePost.GET("/cate/:cid", v1.GetPostByCate)
		}
		// 登录接口
		routeV1Public.POST("login", v1.Login)
	}

	return r
}
