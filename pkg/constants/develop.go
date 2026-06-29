package constants

// HTTP内容类型
const (
	JsonContentType              = "application/json"
	FileContentType              = "application/octet-stream"
	FormUrlEncodedContentType    = "application/x-www-form-urlencoded"
	MultipartFormDataContentType = "multipart/form-data"
)

// 加密算法
const (
	// 单向加密算法（哈希算法）
	MD5    = "MD5"
	SHA1   = "SHA1"
	SHA256 = "SHA256"
	SHA512 = "SHA512"
	SM3    = "SM3"

	// 对称加密算法
	AES      = "AES"
	SM4      = "SM4"
	DES      = "DES"
	ThreeDES = "3DES"
	ChaCha20 = "ChaCha20"
	RC4      = "RC4"

	// 非对象加密算法
	RSA = "RSA"
	ECC = "ECC"
	DSA = "DSA"
	SM2 = "SM2"
)

// 工作模式常量
const (
	ModeECB = "ECB"
	ModeCBC = "CBC"
	ModeGCM = "GCM"
)

// 填充方式常量
const (
	// 对称加密中的块加密填充方式
	PKCS7Padding    = "PKCS7"
	ISO10126Padding = "ISO10126"
	NoPadding       = "NoPadding"
	ZeroPadding     = "ZeroPadding"

	// 非对称加密填充方式
	PKCS1v15 = "PKCS1v15"
	OAEP     = "OAEP"
)

// 编码方式常量
const (
	EncodingBase64 = "base64"
	EncodingHex    = "hex"
)
