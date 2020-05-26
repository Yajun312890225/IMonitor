package router

import (
	"os"

	"iMonitor/handler"
	"iMonitor/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 配置路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())
	r.Use(middleware.Logrus())

	// 配置swagger
	swagURL := ginSwagger.URL(os.Getenv("SWAGGER_URL"))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swagURL))
	r.POST("/login", handler.Login)

	// 可自由配置统一入口，比如/api/v1 版本信息
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthCheckRole())
	{
		v1.POST("query", handler.Query)

		v1.GET("/rolelist", handler.GetRoleList)
		v1.GET("/role/:roleId", handler.GetRole)
		v1.POST("/role", handler.InsertRole)
		v1.PUT("/role", handler.UpdateRole)
		// v1.DELETE("/role/:roleId", handler.DeleteRole)

		v1.GET("/menulist", handler.GetMenuList)
		v1.POST("/menu", handler.InsertMenu)
		v1.PUT("/menu", handler.UpdateMenu)
		v1.DELETE("/menu/:id", handler.DeleteMenu)

		// // 用户登录

		// // 需要登录保护的
		// auth := v1.Group("")
		// auth.Use(middleware.AuthRequired())
		// {
		// 	// User Routing
		// 	auth.POST("user/logout", handler.Logout)
		// }
	}
	return r
}
