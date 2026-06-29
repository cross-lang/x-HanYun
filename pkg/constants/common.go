package constants


// 上下文 Key
const (
	CtxKeyXOrigin    = "X_ORIGIN"          // 请求来源信息 key, 值参考: https://xxx.com
	CtxKeyXReqID     = "X_REQUEST_ID"      // 请求 ID Key, 值参考: d95b060a20b9d22b9df0
	CtxKeyXReqRealIP = "X_REQUEST_REAL_IP" // 请求真实 IP Key
)

// 基本属性
const (
	AppId   = "Hanyun"
	AppName = "汉云"
)

// 匹配模式
const (
	MatchModeExact = "exact" // 精确匹配
	MatchModeFuzzy = "fuzzy" // 模糊匹配
)
