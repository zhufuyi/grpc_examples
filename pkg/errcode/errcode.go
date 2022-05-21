package errcode

// 参考 https://github.com/noChaos1012/tour/tree/master/tag-service/pkg/errcode

import (
	"fmt"

	pb "grpc_examples/pkg/errcode/proto/rpcerrorpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// 自定义业务错误码
	ERROR_LOGIN_DAIL = NewError(20001, "登录失败")

	// grpc内部错误码
	Success          = NewError(0, "ok")
	Fail             = NewError(10000, "内部错误")
	InvalidParams    = NewError(10001, "无效参数")
	Unauthorized     = NewError(10002, "认证错误")
	NotFound         = NewError(10003, "没有找到")
	Unknown          = NewError(10004, "未知")
	DeadlineExceeded = NewError(10005, "超出最后止期限")
	AccessDenied     = NewError(10006, "访问被拒绝")
	LimitExceed      = NewError(10007, "访问限制")
	MethodNotAllowed = NewError(10008, "不支持该方法")
)

type Error struct {
	code int
	msg  string
}

var errorCodes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := errorCodes[code]; ok {
		panic(fmt.Sprintf("code %d 已经存在", code))
	}

	errorCodes[code] = msg

	return &Error{code: code, msg: msg}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) String() string {
	return fmt.Sprintf("code: %d, msg: %s", e.code, e.msg)
}

// ToRPCCode 自定义错误码转换为RPC识别的错误码，避免返回Unknown状态码
func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code

	switch code {
	case Fail.code:
		statusCode = codes.Internal
	case InvalidParams.code:
		statusCode = codes.InvalidArgument
	case Unauthorized.code:
		statusCode = codes.Unauthenticated
	case NotFound.code:
		statusCode = codes.NotFound
	case DeadlineExceeded.code:
		statusCode = codes.DeadlineExceeded
	case AccessDenied.code:
		statusCode = codes.PermissionDenied
	case LimitExceed.code:
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.code:
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}

	return statusCode
}

// ----------------------------------------------------------------------------------

type Status struct {
	*status.Status
}

func FromError(err error) *Status {
	s, _ := status.FromError(err)
	return &Status{s}
}

// ToGRPCStatus 除了原始业务错误码，新增其他说明信息msg，主要给内部客户端
func ToGRPCStatus(err *Error, msg string) *Status {
	s, _ := status.New(ToRPCCode(err.code), msg).WithDetails(&pb.Error{Code: int32(err.code), Message: err.msg})
	return &Status{s}
}

// ToGRPCError 通过Details属性返回错误信息给外部客户端
func ToGRPCError(err *Error) error {
	s, _ := status.New(ToRPCCode(err.code), err.msg).WithDetails(&pb.Error{Code: int32(err.code), Message: err.msg})
	return s.Err()
}
