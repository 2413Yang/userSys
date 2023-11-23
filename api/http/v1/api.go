package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user_system/config"
	"user_system/internal/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Ping 健康连接
func Ping(c *gin.Context) {
	appConfig := config.GetGlobalConf().AppConf
	confInfo, _ := json.MarshalIndent(appConfig, "", "  ")
	appInfo := fmt.Sprintf("app_name: %s\nversion: %s\n\n%s", appConfig.AppName, appConfig.Version, string(confInfo))
	c.String(http.StatusOK, appInfo)
}

// Register 注册
func Register(c *gin.Context) {
	req := &service.RegisterRequest{}
	rsp := &HttpResponse{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("request json err %v", err)
		rsp.ResponseWithError(c, CodeBodyBindErr, err.Error())
		return
	}
	if err := service.Register(req); err != nil {
		rsp.ResponseWithError(c, CodeRegisterErr, err.Error())
		return
	}
	rsp.ResponseSuccess(c)
}
