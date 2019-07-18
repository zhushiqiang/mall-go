package logic

import (
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis"
	"mall-go/common/baseerror"
	"mall-go/user/model"
	"strconv"
)

type (
	UserLogic struct {
		redisCache *redis.Client
		userModel *model.UserModel
	}
	RegisterRequest struct {
		Mobile string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	RegisterResponse struct {

	}
	LoginRequest struct {
		Mobile string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResponse struct {
		Authorization string `json:"authorization"`
	}
)



func NewUserLogic(userModel *model.UserModel, redisCache *redis.Client) *UserLogic {

	return &UserLogic{userModel: userModel, redisCache:redisCache}
}
//注册方法
var ErrorRecordExit = baseerror.NewBaseError("此手机号已经存在")
func (l *UserLogic) Register(r *RegisterRequest) (*RegisterResponse, error) {
	response := new(RegisterResponse)
	b,err  := l.userModel.ExistMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	if b {
		return nil, ErrorRecordExit
	}
	user := &model.User{
		Mobile:r.Mobile,
		Password:fmt.Sprintf("%x", md5.Sum([]byte(r.Password))),
	}
	if _,err := l.userModel.Insert(user); err != nil {
		return nil, err
	}
	return response, nil
}

func (l *UserLogic) Login(r *LoginRequest) (*LoginResponse, error) {
	response := new(LoginResponse)
	user,err  := l.userModel.FindByMobile(r.Mobile)
	if err != nil {
		return nil, err
	}
	response.Authorization = fmt.Sprintf("%x", md5.Sum([]byte(user.Mobile+strconv.Itoa(int(user.Id)))))
	l.redisCache.Set(response.Authorization, user.Id, model.AuthorizationExpire)
	return response, nil
}