package router

import (
	"strconv"
	api "user_system/api/http/v1"
	"user_system/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 路由配置、启动服务
func InitRouterAndServe() {
	setAppRunMode()
	r := gin.Default()

	//健康检查
	r.GET("ping", api.Ping)
	//用户注册
	r.POST("/user/register", api.Register)
	port := config.GetGlobalConf().AppConf.Port
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Error("start server err:" + err.Error())
	}
}

// 设置运行模式
func setAppRunMode() {
	if config.GetGlobalConf().AppConf.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
