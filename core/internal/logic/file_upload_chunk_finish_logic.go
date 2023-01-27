package logic

import (
	"context"
	"path"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadChunkFinishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadChunkFinishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadChunkFinishLogic {
	return &FileUploadChunkFinishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadChunkFinishLogic) FileUploadChunkFinish(req *types.FileUploadChunkFinishRequest) (resp *types.FileUploadChunkFinishReply, err error) {
	co := make([]cos.Object, 0)
	for _, v := range req.CosObjects {
		co = append(co, cos.Object{
			ETag:       v.Etag,
			PartNumber: v.PartNumber,
		})
	}

	err = helper.FinishPartUpload(req.Key, req.UploadId, co)
	if err != nil {
		return
	}

	// 数据入库
	rp := &models.RepositoryPool{
		Identity: helper.UUID(),
		Hash:     req.Hash,
		Name:     req.Key,
		Ext:      path.Ext(req.Key),
		Size:     req.Size,
		Path:     req.Path,
	}
	_, err = l.svcCtx.Engine.Insert(rp)
	if err != nil {
		return
	}

	// 返回结果
	resp = new(types.FileUploadChunkFinishReply)
	resp.Name = req.Key
	resp.Ext = path.Ext(req.Key)
	resp.Identity = rp.Identity

	return
}
