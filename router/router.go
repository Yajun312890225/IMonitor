package router

import (
	"os"

	"iMonitor/handler"
	"iMonitor/middleware"

	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter 配置路由
func InitRouter() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()

	r.Static("/static", "./static")
	handler.Serve()
	r.GET("/ws", handler.ServerLog) //websocket

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())
	r.Use(middleware.Logrus())

	// 配置swagger
	swagURL := ginSwagger.URL(os.Getenv("SWAGGER_URL"))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swagURL))

	// 可自由配置统一入口，比如/api/v1 版本信息
	v1 := r.Group("/api/v1")
	{

		v1.POST("/login", handler.Login)
		v1.POST("/logout", handler.Logout)
		v1.Use(middleware.AuthRequired())
		v1.Use(middleware.AuthCheckRole())

		v1.GET("/query", handler.Query)

		v1.GET("/getinfo", handler.GetInfo)

		v1.GET("/rolelist", handler.GetRoleList)
		v1.GET("/role/:roleId", handler.GetRole)
		v1.POST("/role", handler.InsertRole)
		v1.PUT("/role", handler.UpdateRole)
		v1.DELETE("/role/:roleId", handler.DeleteRole)

		v1.GET("/menulist", handler.GetMenuList)
		v1.GET("/menu/:menuId", handler.GetMenu)
		v1.POST("/menu", handler.InsertMenu)
		v1.PUT("/menu", handler.UpdateMenu)
		v1.DELETE("/menu/:menuId", handler.DeleteMenu)
		v1.GET("/menurole", handler.GetMenuRole)

		v1.GET("/userlist", handler.GetUserList)
		v1.GET("/user/:userId", handler.GetUser)
		v1.GET("/alluser", handler.GetAllUser)
		v1.POST("/user", handler.InsertUser)
		v1.PUT("/user", handler.UpdateUser)
		v1.DELETE("/user/:userId", handler.DeleteUser)
		v1.POST("/user/profileAvatar", handler.InsertUserAvatar)

		server := v1.Group("")
		{

			server.GET("/serverlist", handler.GetServerList)
			server.POST("/ping", handler.Ping)
			server.POST("/server", handler.AddServer)

			//url有参数的请求权限控制写在自己的方法里
			server.GET("/server/:serverId", handler.GetServer)
			server.DELETE("/server/:serverId", handler.DeleteServer)
			server.GET("/querysync/:serverId", handler.QuerySyncOrgId)
			server.POST("/uploadfile", handler.UploadServerFile)
			server.GET("/restartserver/:serverId", handler.RestartServer)
			server.GET("/queryclientlog", handler.QueryClientLog)
			server.GET("/getClientVersion/:serverId", handler.GetClientVersion)

			server.Use(middleware.CheckPermission())
			server.PUT("/server", handler.UpdateServer)
			server.POST("/fetchContact", handler.FetchContact)
			server.POST("/searchUser", handler.SearchUser)
			server.POST("/fetchUserGroup", handler.FetchUserGroup)
			server.POST("/fetchMsgRecord", handler.FetchMsgRecord)
			server.POST("/createcollaborator", handler.CreateCollaborator)
			server.DELETE("/removecollaborator", handler.RemoveCollaborator)
			server.POST("/syncContacts", handler.SyncContacts)
			server.POST("/updatesync", handler.UpdateSync)
			server.POST("/serviceInfo", handler.GetServiceInfo)
			server.POST("/updateVersion", handler.UpdateVersion)

		}

	}
	return r
}
