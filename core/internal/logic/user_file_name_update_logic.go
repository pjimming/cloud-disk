package logic

import (
	"context"
	"errors"

	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest, userIdentity string) (resp *types.UserFileNameUpdateReply, err error) {
	// 判断当前名称在该层级是否存在
	cnt, err := l.svcCtx.Engine.
		Where("name = ? AND parent_id = (SELECT parent_id FROM user_repository ur WHERE ur.identity = ?)", req.Name, req.Identity).
		Count(new(models.UserRepository))
	if err != nil {
		return
	}
	if cnt > 0 {
		return nil, errors.New("文件名称已存在")
	}

	// 用户存储池文件名称修改
	data := &models.UserRepository{Name: req.Name}
	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Update(data)
	if err != nil {
		return
	}
	// 公共存储池文件名称修改
	ur := new(models.UserRepository)
	_, err = l.svcCtx.Engine.Where("identity = ? AND user_identity = ?", req.Identity, userIdentity).Get(ur)
	if err != nil {
		return
	}
	rp := &models.RepositoryPool{Name: req.Name}
	_, err = l.svcCtx.Engine.Where("identity = ?", ur.RepositoryIdentity).Update(rp)
	return
}
