package response

// Res 请求返回统一格式
type Res struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

type Page struct {
	List      interface{} `json:"list"`
	Count     int         `json:"count"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
}

type PageResponse struct {
	Code int    `json:"code" example:"200"`
	Data Page   `json:"data"`
	Msg  string `json:"msg"`
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeSuccess 成功
	CodeSuccess = 0
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoRightErr 未授权访问
	CodeNoRightErr = 403
	// CodeDBError 数据库操作失败
	CodeDBError = 50001
	// CodeEncryptError 加密失败
	CodeEncryptError = 50002
	// CodeParamErr 各种奇奇怪怪的参数错误
	CodeParamErr = 40001
	// CodeUserNotFound 用户不存在
	CodeUserNotFound = 40002
	// CodePasswordErr 密码错误
	CodePasswordErr = 40003

	//CodeAccessionNotPermission 没有权限
	CodeAccessionNotPermission = 40004
)

// CodeErrMsg 服务器错误码对应错误信息
var CodeErrMsg = map[int]string{
	CodeCheckLogin:             "未登录",
	CodeNoRightErr:             "未授权访问",
	CodeDBError:                "数据库操作失败",
	CodeEncryptError:           "加密失败",
	CodeParamErr:               "参数错误",
	CodeUserNotFound:           "用户不存在",
	CodePasswordErr:            "密码错误",
	CodeAccessionNotPermission: "没有权限",
}
