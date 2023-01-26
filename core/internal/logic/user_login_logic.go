package logic

import (
	"context"
	"errors"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.UserLoginReply, err error) {
	// 1. 从数据库里面查询当前用户
	user := new(models.UserBasic)
	has, err := l.svcCtx.Engine.Where("name = ? AND password = ?", req.Name, helper.Md5(req.Password)).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户名或密码错误！")
	}
	// 2. 生成token
	token, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.TokenExpireTime)
	if err != nil {
		return nil, err
	}
	// 3. 生成refreshToken
	refreshToken, err := helper.GenerateToken(user.Id, user.Identity, user.Name, define.RefreshTokenExpireTime)
	if err != nil {
		return nil, err
	}
	// 4. 返回结果
	resp = new(types.UserLoginReply)
	resp.Token = token
	resp.RefreshToken = refreshToken
	return
}
