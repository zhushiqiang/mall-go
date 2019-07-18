package logic

import (
	"context"
	"mall-go/user/model"
	"mall-go/user/protos"
)

type (
	UserRpcServerLogic struct {
		userModel *model.UserModel
	}
)

//	FindUserByMobile(context.Context, *FindUserByMobileRequest) (*UserResponse, error)
//	FindId(context.Context, *FindIdRequest) (*UserResponse, error)

func NewUserRpcServiceLogic (userModel *model.UserModel) *UserRpcServerLogic {
	return &UserRpcServerLogic{userModel:userModel}
}

func (l *UserRpcServerLogic) FindUserByMobile(_ context.Context, r *protos.FindUserByMobileRequest) (*protos.UserResponse, error) {
	user, err := l.userModel.FindByMobile(r.Mobile)

	if err != nil {
		return nil,err
	}

	return &protos.UserResponse{
		Id:user.Id,Name:user.Name,
	},nil
}

func (l *UserRpcServerLogic) FindId(_ context.Context, r *protos.FindIdRequest) (*protos.UserResponse, error) {
	user, err := l.userModel.FindById(r.Id)

	if err != nil {
		return nil,err
	}

	return &protos.UserResponse{
		Id:user.Id,Name:user.Name,
	},nil
}
