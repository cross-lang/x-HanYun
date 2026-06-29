package oplog

import (
	"net/http"

	"x-HanYun/internal/logic/oplog"
	"x-HanYun/internal/svc"
	"x-HanYun/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateOpLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOpLogRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		logx.WithContext(r.Context()).Info("===========================")
		logx.WithContext(r.Context()).Infof(">>>>>>>> CreateOpLog httpx.Parse req: %+v", req)

		l := oplog.NewCreateOpLogLogic(r.Context(), svcCtx)
		resp, err := l.CreateOpLog(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
