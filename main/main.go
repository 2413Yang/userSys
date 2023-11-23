package main

import (
	"user_system/config"
	"user_system/internal/router"
	"user_system/utils"
)

func main() {
	config.InitConfig()
	utils.GetDB()
	router.InitRouterAndServe()
}
