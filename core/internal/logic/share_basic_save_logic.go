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

type ShareBasicSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicSaveLogic {
	return &ShareBasicSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicSaveLogic) ShareBasicSave(req *types.ShareBasicSaveRequest, userIdentity string) (resp *types.ShareBasicSaveReply, err error) {
	// 找到分享的记录
	sb := new(models.ShareBasic)
	has, err := l.svcCtx.Engine.Where("identity = ?", req.Identity).Get(sb)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("该分享记录不存在")
	}
	// 从 repository_pool 内获取资源信息
	rp := new(models.RepositoryPool)
	has, err = l.svcCtx.Engine.Where("identity = ?", sb.RepositoryIdentity).Get(rp)
	if err != nil {
		return
	}
	if !has {
		return nil, errors.New("资源不存在")
	}
	// 把资源保存到 user_repository 里面
	ur := &models.UserRepository{
		Identity:           helper.UUID(),
		UserIdentity:       userIdentity,
		ParentId:           req.ParentId,
		Ext:                rp.Ext,
		RepositoryIdentity: rp.Identity,
		Name:               rp.Name,
	}
	_, err = l.svcCtx.Engine.Insert(ur)
	// 返回结果
	resp = new(types.ShareBasicSaveReply)
	resp.Identity = ur.Identity
	return
}
