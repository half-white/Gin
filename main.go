package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// 自定义一个中间件 拦截器
func myHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		//通过自定义的中间件，设置的值，在后续处理只要调用了这个中间件的都可以拿到这里的参数
		context.Set("usersession", "userid-1")
		context.Next() //放行
		//context.Abort()//拦截
	}
}

func main() {
	//创建一个服务
	ginServer := gin.Default()
	//注册网页头像
	ginServer.Use(favicon.New("./favicon.ico"))

	//注册中间件
	ginServer.Use(myHandler())

	//连接数据库的代码

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")
	//加载资源文件
	//ginServer.Static("/static", "./static")

	//响应一个页面给前端
	ginServer.GET("/index", func(context *gin.Context) {
		//context.JSON() json 数据
		context.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "这是go后台传递的数据",
		})
	})

	//接收前端传递的参数
	// usl? userid=xxx&username=xep  传统参数
	ginServer.GET("/user/info", myHandler(), func(context *gin.Context) {

		//取出中间件里的值
		usersession := context.MustGet("usersession").(string)
		log.Println("============", usersession)

		userid := context.Query("userid")
		username := context.Query("username")
		context.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})

	})

	// /user/info/1/xep              RESTful参数
	ginServer.GET("/user/info/:userid/:username", func(context *gin.Context) {
		userid := context.Param("userid")
		username := context.Param("username")
		context.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	//前端给后端传递JSON
	ginServer.POST("/json", func(context *gin.Context) {
		//request.body
		data, _ := context.GetRawData()
		var m map[string]interface{}
		//包装为json数据 []byte
		_ = json.Unmarshal(data, &m)
		context.JSON(http.StatusOK, m)
	})

	//访问地址 处理请求
	ginServer.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "hello word"})
	})

	ginServer.POST("/user/add", func(context *gin.Context) {
		username := context.PostForm("username")
		password := context.PostForm("password")

		//可以增加数据校验
		context.JSON(http.StatusOK, gin.H{
			"msg":      "ok",
			"username": username,
			"password": password,
		})
	})

	//路由
	ginServer.GET("/test", func(context *gin.Context) {
		//重定向
		context.Redirect(http.StatusMovedPermanently, "https://www.kuangstudy.com")
	})

	//404
	ginServer.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})

	//路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
		userGroup.GET("/login")
		userGroup.GET("/logout")
	}
	orderGroup := ginServer.Group("/order")
	{
		orderGroup.GET("/add")
		orderGroup.DELETE("/delete")
	}

	//中间件

	//服务器端口
	ginServer.Run(":8082")

}
