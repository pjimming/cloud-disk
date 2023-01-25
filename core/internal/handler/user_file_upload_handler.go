package handler

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"path"

	"cloud-disk/core/helper"
	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"cloud-disk/core/models"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserFileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// start handle
		// 判断文件是否存在
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			return
		}
		b := make([]byte, fileHeader.Size)
		_, err = file.Read(b) // 通过reader写入文件
		if err != nil {
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum(b)) // 计算文件的md5值，获取hash值
		rp := new(models.RepositoryPool)
		has, err := svcCtx.Engine.Where("hash = ?", hash).Get(rp) // 从rp数据库中查找数据
		if err != nil {
			return
		}
		if has { // 有数据，返回相应信息
			httpx.OkJson(w, &types.UserFileUploadReply{Identity: rp.Identity, Ext: rp.Ext, Name: rp.Name})
			return
		}
		// 上传文件到COS中
		cosPath, err := helper.CosUpload(r)
		if err != nil {
			return
		}

		// 给logic传递数据
		req.Name = fileHeader.Filename
		req.Ext = path.Ext(fileHeader.Filename)
		req.Size = fileHeader.Size
		req.Path = cosPath
		req.Hash = hash
		// end handle

		l := logic.NewUserFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.UserFileUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
