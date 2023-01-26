package logic

import (
	"context"

	"cloud-disk/core/define"
	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshAuthorizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAuthorizationLogic {
	return &RefreshAuthorizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshAuthorizationLogic) RefreshAuthorization(req *types.RefreshAuthorizationRequest, auth string) (resp *types.RefreshAuthorizationReply, err error) {
	// 解析token，获取uc
	uc, err := helper.AnalyzeToken(auth)
	if err != nil {
		return
	}
	// 生成token
	token, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.TokenExpireTime)
	if err != nil {
		return
	}
	// 生成refreshtoken
	refreshToken, err := helper.GenerateToken(uc.Id, uc.Identity, uc.Name, define.RefreshTokenExpireTime)
	if err != nil {
		return
	}
	// 返回结果
	resp = new(types.RefreshAuthorizationReply)
	resp.Token = token
	resp.RefreshToken = refreshToken

	return
}
