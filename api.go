package log

import "os"

var (
	// 默认
	defaultLogger Logger
	// 级别
	debugLevel = "[debug] "
	infoLevel  = "[info] "
	warnLevel  = "[warn] "
	errorLevel = "[error] "
	panicLevel = "[panic] "
)

func init() {
	// 默认 logger
	defaultLogger = NewLogger(os.Stdout, new(DefaultHeader), "")
}

// SetLogger 设置默认 Logger
func SetLogger(lg Logger) {
	defaultLogger = lg
}

// GetLogger 获取默认 Logger
func GetLogger() Logger {
	return defaultLogger
}

// Recover 使用默认 Logger 函数
func Recover(recover any) {
	defaultLogger.Recover(recover)
}

// IsDebug 使用默认 Logger 函数
func IsDebug() bool {
	return defaultLogger.IsDebug()
}

// EnableDebug 使用默认 Logger 函数
func EnableDebug(enable bool) {
	defaultLogger.EnableDebug(enable)
}

// Debug 使用默认 Logger 函数
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Debugf 使用默认 Logger 函数
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// DebugDepth 使用默认 Logger 函数
func DebugDepth(depth int, args ...any) {
	defaultLogger.DebugDepth(depth, args...)
}

// DebugfDepth 使用默认 Logger 函数
func DebugfDepth(depth int, format string, args ...any) {
	defaultLogger.DebugfDepth(depth, format, args...)
}

// DebugTrace 使用默认 Logger 函数
func DebugTrace(traceID string, args ...any) {
	defaultLogger.DebugTrace(traceID, args...)
}

// DebugfTrace 使用默认 Logger 函数
func DebugfTrace(traceID string, format string, args ...any) {
	defaultLogger.DebugfTrace(traceID, format, args...)
}

// DebugDepthTrace 使用默认 Logger 函数
func DebugDepthTrace(depth int, traceID string, args ...any) {
	defaultLogger.DebugDepthTrace(depth, traceID, args...)
}

// DebugfDepthTrace 使用默认 Logger 函数
func DebugfDepthTrace(depth int, traceID string, format string, args ...any) {
	defaultLogger.DebugfDepthTrace(depth, traceID, format, args...)
}

// IsInfo 使用默认 Logger 函数
func IsInfo() bool {
	return defaultLogger.IsInfo()
}

// EnableInfo 使用默认 Logger 函数
func EnableInfo(enable bool) {
	defaultLogger.EnableInfo(enable)
}

// Info 使用默认 Logger 函数
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Infof 使用默认 Logger 函数
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// InfoDepth 使用默认 Logger 函数
func InfoDepth(depth int, args ...any) {
	defaultLogger.InfoDepth(depth, args...)
}

// InfofDepth 使用默认 Logger 函数
func InfofDepth(depth int, format string, args ...any) {
	defaultLogger.InfofDepth(depth, format, args...)
}

// InfoTrace 使用默认 Logger 函数
func InfoTrace(traceID string, args ...any) {
	defaultLogger.InfoTrace(traceID, args...)
}

// InfofTrace 使用默认 Logger 函数
func InfofTrace(traceID string, format string, args ...any) {
	defaultLogger.InfofTrace(traceID, format, args...)
}

// InfoDepthTrace 使用默认 Logger 函数
func InfoDepthTrace(depth int, traceID string, args ...any) {
	defaultLogger.InfoDepthTrace(depth, traceID, args...)
}

// InfofDepthTrace 使用默认 Logger 函数
func InfofDepthTrace(depth int, traceID string, format string, args ...any) {
	defaultLogger.InfofDepthTrace(depth, traceID, format, args...)
}

// IsWarn 使用默认 Logger 函数
func IsWarn() bool {
	return defaultLogger.IsWarn()
}

// EnableWarn 使用默认 Logger 函数
func EnableWarn(enable bool) {
	defaultLogger.EnableWarn(enable)
}

// Warn 使用默认 Logger 函数
func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

// Warnf 使用默认 Logger 函数
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// WarnDepth 使用默认 Logger 函数
func WarnDepth(depth int, args ...any) {
	defaultLogger.WarnDepth(depth, args...)
}

// WarnfDepth 使用默认 Logger 函数
func WarnfDepth(depth int, format string, args ...any) {
	defaultLogger.WarnfDepth(depth, format, args...)
}

// WarnTrace 使用默认 Logger 函数
func WarnTrace(traceID string, args ...any) {
	defaultLogger.WarnTrace(traceID, args...)
}

// WarnfTrace 使用默认 Logger 函数
func WarnfTrace(traceID string, format string, args ...any) {
	defaultLogger.WarnfTrace(traceID, format, args...)
}

// WarnDepthTrace 使用默认 Logger 函数
func WarnDepthTrace(depth int, traceID string, args ...any) {
	defaultLogger.WarnDepthTrace(depth, traceID, args...)
}

// WarnfDepthTrace 使用默认 Logger 函数
func WarnfDepthTrace(depth int, traceID string, format string, args ...any) {
	defaultLogger.WarnfDepthTrace(depth, traceID, format, args...)
}

// IsError 使用默认 Logger 函数
func IsError() bool {
	return defaultLogger.IsError()
}

// EnableError 使用默认 Logger 函数
func EnableError(enable bool) {
	defaultLogger.EnableError(enable)
}

// Error 使用默认 Logger 函数
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Errorf 使用默认 Logger 函数
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// ErrorDepth 使用默认 Logger 函数
func ErrorDepth(depth int, args ...any) {
	defaultLogger.ErrorDepth(depth, args...)
}

// ErrorfDepth 使用默认 Logger 函数
func ErrorfDepth(depth int, format string, args ...any) {
	defaultLogger.ErrorfDepth(depth, format, args...)
}

// ErrorTrace 使用默认 Logger 函数
func ErrorTrace(traceID string, args ...any) {
	defaultLogger.ErrorTrace(traceID, args...)
}

// ErrorfTrace 使用默认 Logger 函数
func ErrorfTrace(traceID string, format string, args ...any) {
	defaultLogger.ErrorfTrace(traceID, format, args...)
}

// ErrorDepthTrace 使用默认 Logger 函数
func ErrorDepthTrace(depth int, traceID string, args ...any) {
	defaultLogger.ErrorDepthTrace(depth, traceID, args...)
}

// ErrorfDepthTrace 使用默认 Logger 函数
func ErrorfDepthTrace(depth int, traceID string, format string, args ...any) {
	defaultLogger.ErrorfDepthTrace(depth, traceID, format, args...)
}
