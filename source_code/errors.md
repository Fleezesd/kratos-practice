// errors.go
```go
package errors

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is a status error.
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.Code, e.Reason, e.Message, e.Metadata, e.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.Metadata = md
	return err
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	return &Error{
		Status: Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
		},
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...interface{}) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: metadata,
		},
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, UnknownReason, err.Error())
	}
	ret := New(
		httpstatus.FromGRPCCode(gs.Code()),
		UnknownReason,
		gs.Message(),
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.Reason = d.Reason
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

```

// kratos error 简易介绍
```go
Kratos定义了错误码和错误信息的规范,在errors/errors.go中:
错误码为int类型,推荐使用5位数字编码
错误信息为字符串,支持多语言翻译
通过New错误构造函数创建规范化的错误
errors包实现了错误码的定义、错误信息翻译等功能。

httputil/http.go实现了HTTP错误的规范化处理,可以把Kratos错误翻译成标准HTTP错误返回。

grpc/grpc.go实现了gRPC错误翻译,可以把Kratos错误翻译成gRPC状态码和错误详情。

middleware/recovery可以捕获Panic作为错误返回。

validator可以在请求入参中捕获验证错误。

logger日志可以捕获并记录错误堆栈。

trace中间件可以追踪错误并上报。

通过错误码可以对错误进行分类处理。
```

// http & grpc status
```go
- 错误码(Code)
HTTP使用状态码表示,如400,500等
gRPC使用grpc状态码,如Unknown、InvalidArgument等
- 错误原因(Reason)
字符串,表示错误的具体原因
- 错误信息(Message)
字符串,表示错误详情描述
- 元数据(Metadata)
键值对,用于传递额外调试信息
- 原始错误(Cause)
记录原始错误对象

翻译机制
Kratos errors包实现了HTTP状态码和gRPC状态码的相互转换
传递机制
HTTP通过状态码、错误信息、元数据头
gRPC通过grpc状态和错误详情
所以Kratos通过errors实现了错误在HTTP和gRPC之间的自然转换,保证了服务间通用的错误处理格式,方便错误跟踪和监控。

总体来说,核心是错误码、错误原因、错误信息这些通用字段,以及翻译和传递机制,来实现HTTP和gRPC的错误统一
```


// errors.go 实现方法
```go
New/Newf方法,用于创建和格式化错误

Errorf方法,把错误包装为error接口

Code/Reason方法,获取错误码和原因

FromError方法,把普通error转换为Error对象

Clone方法,深拷贝错误对象

WithXXX方法,设置错误的原因、消息和元数据
```

// errors.go 赋值with时采用clone方式 原因
```go
- 避免共享内存,减少副作用
错误对象可能在不同的地方被多次使用,如果直接修改会产生副作用,所以需要克隆一个新的错误对象来设置额外信息。

- 保证错误对象的不变性
错误一旦创建,其信息就不应该再被修改,以保证错误状态的一致性。克隆可以避免修改原错误对象。

- 扩展错误信息
通过克隆原错误并设置新的元信息,可以在不改变原错误的基础上扩展新的调试信息。

- 遵循 Go 错误处理最佳实践
Go 错误处理的最佳实践是不要修改共享的错误对象,克隆一个新的错误对象来添加信息。

- 保持向下兼容
克隆操作不会破坏之前的错误处理逻辑,可以平滑升级。

所以 Kratos 的错误处理通过 Clone 来获取一个独立的错误对象副本,在此基础上进行扩展,可以避免共享内存的副作用,保证错误对象的不变性,也使错误处理更符合 Go 的最佳实践。
```