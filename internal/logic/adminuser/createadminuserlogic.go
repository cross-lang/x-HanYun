package adminuser

import (
	"context"

	"x-HanYun/internal/svc"
	"x-HanYun/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminUserLogic {
	return &CreateAdminUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAdminUserLogic) CreateAdminUser(req *types.CreateAdminUserRequest) (resp *types.CreateAdminUserResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
