package xerrors

import (
	"context"
	xerrors "errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"go.uber.org/zap"

	"x-HanYun/pkg/constants"
	"x-HanYun/pkg/core/log"
)

// HTTPErrorHandler 处理错误并返回HTTP状态码和响应数据
// 参数 ctx 是上下文
// 参数 err 是错误信息
// 返回值是HTTP状态码和响应数据
func HTTPErrorHandler(ctx context.Context, err error) (int, interface{}) {
	var codeError Error
	switch {
	case xerrors.As(err, &codeError):
		// 处理业务错误
		httpCode := mapErrorCodeToHTTP(codeError.Code())
		log.WithContext(ctx).Error("业务错误",
			zap.Int("code", int(codeError.Code())),
			zap.String("msg", codeError.Message()),
			zap.String("debug", codeError.DebugInfo()))

		return httpCode, constants.BaseResponse{
			Code:  int(codeError.Code()),
			Msg:   codeError.Message(),
			Debug: fmt.Sprintf("traceId:%s, %s", trace.TraceIDFromContext(ctx), codeError.DebugInfo()),
		}

	default:
		// 处理系统错误
		log.WithContext(ctx).Error("系统错误", zap.Error(err))
		return http.StatusInternalServerError, constants.BaseResponse{
			Code:  -1,
			Msg:   "系统内部错误",
			Debug: fmt.Sprintf("traceId:%s, error:%v", trace.TraceIDFromContext(ctx), err),
		}
	}
}
