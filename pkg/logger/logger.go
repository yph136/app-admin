package logger

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// 定义建名
const (
	// 开始时间
	StartedAtKey = "started_at"
	// 追踪 ID
	TraceIDKey = "trace_id"
	// 用户 ID
	UserIDKey = "user_id"
	// 接口 ID
	SpanIDKey = "span_id"
	// 接口主题
	SpanTitleKey = "span_title"
	// 函数名
	SpanFunctionKey = "span_function"
	// 版本号
	VersionKey = "version"
	// 校验时间
	TimeConsumingKey = "time_consuming"
)

type TraceIDFunc func() string

var (
	version     string
	TraceIDFunc TraceIDFunc
)

type Logger = logrus.Logger

type Hook = logrus.Hook

// StandardLogger
func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

// SetLevel
func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

//  SetFormatter
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// SetOutput 设定日志输出
func SetOutput(out io.Writer) {
	logrus.Setoutput(out)
}

// SetVersion 设定版本
func SetVersion(v string) {
	version = v
}

// SetTraceIDFunc 设定追踪 ID 的处理函数
func SetTraceIDFunc(fn TraceIDFunc) {
	traceIDFunc = fn
}

// AddHook
func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

// getTraceID
func getTraceID() string {
	if traceIDFunc != nil {
		return traceIDFunc()
	}
	return time.Now().Format("2006.01.02.15.04.05.000")
}

type (
	traceIDContextKey struct{}
	spanIDContextKey  struct{}
	userIDContextKey  struct{}
)

// NewTraceIDContext 创建追踪 ID 上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取追踪 ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID{}
}

// NewSpanIDContext 创建跟踪单元ID
func NewTraceIDContext(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDContextKey{}, spanID)
}

// FromSpanIDContext 从上下文中获取跟踪单元 ID
func FromSpanIDContext(ctx context.Context) string {
	v := ctx.Value(spanIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return getTraceID()
}

// NewUserIDContext 创建用户 ID 上下文
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

// FromUserIDContext 从上下文中获取用户 ID
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// StartSpan 开始一个追踪单元
func StartSpan(ctx context.Context, title, funcName string) *Entry {
	fileds := map[string]interface{}{
		StartedAtKey:   time.Now(),
		UserIDKey:      FromUserIDContext(ctx),
		TraceIDKey:     FromTraceIDContext(ctx),
		SpanIDKey:      FromSpanIDContext(ctx),
		SpanTitleKey:   title,
		SpanFuncionKey: funcName,
		VersionKey:     version,
	}

	return newEntry(logrus.WithFields(fileds))
}

//  StartSpanWithCall 开始一个追踪单元(回调执行)
func StartSpanWithCall(ctx context.Context, title, funcName string) func() *Entry {
	return func() *Entry {
		return StartSpan(ctx, title, funcName)
	}
}

// newEntry
func newEntry(entry *logrus.Entry) *Entry {
	return &Entry{entry: entry}
}

// Entry 定义统一的日志写入方式
type Entry struct {
	entry  *logrus.Entry
	finish int32
}

// Finish 完成, 如果没有触发写入则手动触发  Info 级别的日志写入
func (e *Entry) Finish() {
	if atomic.CompareAndSwapInt32(&e.finish, 0, 1) {
		e.done()
		e.entry.Info()
	}
}

func (e *Entry) checkAndDelete(fileds map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if _, ok := fileds[key]; ok {
			delete(fileds, key)
		}
	}
}

// WithFields 结构化字段写入
func (e *Entry) WithFields(fileds map[string]interface{}) *Entry {
	e.checkAndDelete(fileds,
		StartAtKey,
		TraceIDKey)
}
