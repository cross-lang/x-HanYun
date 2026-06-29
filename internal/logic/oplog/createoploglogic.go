package oplog

import (
	"context"

	"x-HanYun/internal/svc"
	"x-HanYun/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOpLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOpLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOpLogLogic {
	return &CreateOpLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOpLogLogic) CreateOpLog(req *types.CreateOpLogRequest) (resp *types.CreateOpLogResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
