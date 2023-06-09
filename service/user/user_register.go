package user

import (
	"TikTok_Project/repository"
	"TikTok_Project/utils"
	"errors"
)

func UserRegister(name string, password string) (*LoginAndRegisterResponse, int, error) {

	if err := IsValidUser(name, password); err != nil {
		return nil, 1, err
	}

	// 判断用户名是否存在
	userExistDao := repository.InitUserDao()
	if userExistDao.IsExistName(name) {
		return nil, 2, errors.New("用户名已存在")
	}

	var user repository.User
	user.Username = name
	user.Password = password
	user.Avatar = "https://blog-1314857283.cos.ap-shanghai.myqcloud.com/background-img/avatar.jpg"
	user.BackgroundImage = "https://blog-1314857283.cos.ap-shanghai.myqcloud.com/images/202304071526669.jpg"
	user.Signature = "Talk is cheap,show me the code!"

	// 密码加密
	salt := utils.GenerateSalt(16)
	user.Password = utils.HashPassword(password, salt) + ":" + salt

	// 数据库更新用户数据
	userUpdateDao := repository.InitUserDao()
	err := userUpdateDao.UserRegister(&user)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{"status_code": 3, "status_msg": err.Error()})
		return nil, 3, err
	}

	// 获取 token
	token, err := utils.GenToken(user)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{"status_code": 4, "status_msg": err.Error()})
		return nil, 4, err
	}

	response := &LoginAndRegisterResponse{
		UserId: user.ID,
		Token:  token,
	}
	return response, 0, nil
}
