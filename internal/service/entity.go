package service

import (
	"fmt"
	"user_system/internal/dao"
	"user_system/internal/model"
	"user_system/pkg/constant"
	"user_system/utils"

	log "github.com/sirupsen/logrus"
)

// Register 用户注册
func Register(req *RegisterRequest) error {
	if req.UserName == "" || req.Password == "" || req.Age <= 0 || !utils.Contains([]string{constant.GenderFeMale, constant.GenderMale}, req.Gender) {
		log.Error("register param invalid")
		return fmt.Errorf("register param  invalid")
	}
	existedUser, err := dao.GetUserByName(req.UserName)
	if err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("register|%v", err)
	}
	if existedUser != nil {
		log.Errorf("用户已经注册,user_name==%s", req.UserName)
		return fmt.Errorf("用户已经注册,不能重复注册")
	}
	user := &model.User{
		Name:     req.UserName,
		Age:      req.Age,
		Gender:   req.Gender,
		PassWord: req.Password,
		NickName: req.NickName,
		CreateModel: model.CreateModel{
			Creator: req.UserName,
		},
		ModifyModel: model.ModifyModel{
			Modifier: req.UserName,
		},
	}
	log.Infof("user ===== %+v", user)
	if err := dao.CreateUser(user); err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("register|%v", err)
	}
	return nil
}
