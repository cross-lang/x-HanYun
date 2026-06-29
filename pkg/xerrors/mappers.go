package xerrors

import "net/http"

// mapErrorCodeToHTTP 将自定义错误码映射到HTTP状态码
func mapErrorCodeToHTTP(code ErrorCode) int {
	c := int(code)
	switch {
	case c >= 40000000 && c <= 40099999:
		return http.StatusBadRequest
	case c >= 40100000 && c <= 40199999:
		return http.StatusUnauthorized
	case c >= 40300000 && c <= 40399999:
		return http.StatusForbidden
	case c >= 50000000 && c <= 50099999:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
