package logic

import (
	"context"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileUploadLogic {
	return &UserFileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileUploadLogic) UserFileUpload(req *types.UserFileUploadRequest) (resp *types.UserFileUploadReply, err error) {
	// 将handler里面处理好的信息存入数据库
	rp := &models.RepositoryPool{ // 新建rp信息
		Identity: helper.UUID(),
		Hash:     req.Hash,
		Name:     req.Name,
		Ext:      req.Ext,
		Size:     req.Size,
		Path:     req.Path,
	}
	l.svcCtx.Engine.Insert(rp)
	if err != nil {
		return nil, err
	}
	// 返回上传成功的信息
	resp = new(types.UserFileUploadReply)
	resp.Identity = rp.Identity
	resp.Ext = rp.Ext
	resp.Name = rp.Name
	return
}
