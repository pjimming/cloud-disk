package logic

import (
	"context"
	"errors"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicCreateLogic {
	return &ShareBasicCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicCreateLogic) ShareBasicCreate(req *types.ShareBasicCreateRequest, userIdentity string) (resp *types.ShareBasicCreateReply, err error) {
	// 判断文件是否存在
	ur := new(models.UserRepository)
	has, err := l.svcCtx.Engine.Where("user_identity = ? AND repository_identity = ?", userIdentity, req.RepositoryIdentity).Get(ur)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("该文件不存在")
	}
	// 创建分享记录
	data := &models.ShareBasic{
		Identity:           helper.UUID(),
		UserIdentity:       userIdentity,
		RepositoryIdentity: req.RepositoryIdentity,
		ExpiredTime:        req.ExpiredTime,
	}
	_, err = l.svcCtx.Engine.Insert(data)
	if err != nil {
		return
	}
	resp = new(types.ShareBasicCreateReply)
	resp.Identity = data.Identity
	return
}
