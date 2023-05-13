package log

import "os"

var (
	// 默认
	defaultLogger Logger
	// 级别
	debugLevel = "[debug] "
	infoLevel  = "[info] "
	warngLevel = "[warn] "
	errorLevel = "[error] "
	panicLevel = "[panic] "
)

func init() {
	// 默认 logger
	defaultLogger = NewLogger(os.Stdout, DefaultHeader, "")
}

// SetLogger 设置默认 Logger
func SetLogger(lg Logger) {
	defaultLogger = lg
}

// import "io"

// // NewLogger 返回一个 Logger 实例。output 和 headerFormater 是初始化参数。
// func NewLogger(output io.Writer, headerFormater HeaderFormater) Logger {
// 	lg := new(logger)
// 	lg.Writer = output
// 	lg.Header = headerFormater
// 	lg.enableDebug = true
// 	lg.enableInfo = true
// 	lg.enableWarn = true
// 	lg.enableError = true
// 	return lg
// }

// // SetOutput 设置默认 Logger 的输出。
// func SetOutput(output io.Writer) {
// 	defaultLogger.SetOutput(output)
// }

// // SetHeaderFormater 设置默认 Logger 的 HeaderFormater 。
// func SetHeaderFormater(headerFormater HeaderFormater) {
// 	defaultLogger.SetHeaderFormater(headerFormater)
// }

// // SetLevel 设置默认 Logger 的输出级别 。
// func SetLevel(levels ...Level) {
// 	defaultLogger.SetLevel(levels...)
// }

// // IsDebug 返回默认 Logger 是否金鱼 debug 级别。
// func IsDebug() bool {
// 	return defaultLogger.enableDebug
// }

// // Debug 使用默认 Logger 输出日志。
// func Debug(args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.output("", defaultLoggerDepth, DebugLevel, args...)
// 	}
// }

// // Debugf 使用默认 Logger 输出日志。
// func Debugf(format string, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.outputf("", defaultLoggerDepth, DebugLevel, format, args...)
// 	}
// }

// // DebugDepth 使用默认 Logger 输出日志。
// func DebugDepth(depth int, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.output("", depth+defaultLoggerDepth+defaultLoggerDepth, DebugLevel, args...)
// 	}
// }

// // DebugfDepth 使用默认 Logger 输出日志。
// func DebugfDepth(depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.outputf("", depth+defaultLoggerDepth+defaultLoggerDepth, DebugLevel, format, args...)
// 	}
// }

// // DebugTrace 使用默认 Logger 输出日志。
// func DebugTrace(traceID string, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.output(traceID, defaultLoggerDepth, DebugLevel, args...)
// 	}
// }

// // DebugfTrace 使用默认 Logger 输出日志。
// func DebugfTrace(traceID string, format string, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.outputf(traceID, defaultLoggerDepth, DebugLevel, format, args...)
// 	}
// }

// // DebugDepthTrace 使用默认 Logger 输出日志。
// func DebugDepthTrace(traceID string, depth int, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.output(traceID, depth+defaultLoggerDepth+defaultLoggerDepth, DebugLevel, args...)
// 	}
// }

// // DebugfDepthTrace 使用默认 Logger 输出日志。
// func DebugfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableDebug {
// 		defaultLogger.outputf(traceID, depth+defaultLoggerDepth+defaultLoggerDepth, DebugLevel, format, args...)
// 	}
// }

// // IsInfo 返回默认 Logger 是否金鱼 Info 级别。
// func IsInfo() bool {
// 	return defaultLogger.enableInfo
// }

// // Info 使用默认 Logger 输出日志。
// func Info(args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.output("", defaultLoggerDepth, InfoLevel, args...)
// 	}
// }

// // Infof 使用默认 Logger 输出日志。
// func Infof(format string, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.outputf("", defaultLoggerDepth, InfoLevel, format, args...)
// 	}
// }

// // InfoDepth 使用默认 Logger 输出日志。
// func InfoDepth(depth int, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.output("", depth+defaultLoggerDepth+defaultLoggerDepth, InfoLevel, args...)
// 	}
// }

// // InfofDepth 使用默认 Logger 输出日志。
// func InfofDepth(depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.outputf("", depth+defaultLoggerDepth+defaultLoggerDepth, InfoLevel, format, args...)
// 	}
// }

// // InfoTrace 使用默认 Logger 输出日志。
// func InfoTrace(traceID string, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.output(traceID, defaultLoggerDepth, InfoLevel, args...)
// 	}
// }

// // InfofTrace 使用默认 Logger 输出日志。
// func InfofTrace(traceID string, format string, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.outputf(traceID, defaultLoggerDepth, InfoLevel, format, args...)
// 	}
// }

// // InfoDepthTrace 使用默认 Logger 输出日志。
// func InfoDepthTrace(traceID string, depth int, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.output(traceID, depth+defaultLoggerDepth, InfoLevel, args...)
// 	}
// }

// // InfofDepthTrace 使用默认 Logger 输出日志。
// func InfofDepthTrace(traceID string, depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableInfo {
// 		defaultLogger.outputf(traceID, depth+defaultLoggerDepth, InfoLevel, format, args...)
// 	}
// }

// // IsWarn 返回默认 Logger 是否金鱼 Warn 级别。
// func IsWarn() bool {
// 	return defaultLogger.enableWarn
// }

// // Warn 使用默认 Logger 输出日志。
// func Warn(args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.output("", defaultLoggerDepth, WarnLevel, args...)
// 	}
// }

// // Warnf 使用默认 Logger 输出日志。
// func Warnf(format string, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.outputf("", defaultLoggerDepth, WarnLevel, format, args...)
// 	}
// }

// // WarnDepth 使用默认 Logger 输出日志。
// func WarnDepth(depth int, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.output("", depth+defaultLoggerDepth, WarnLevel, args...)
// 	}
// }

// // WarnfDepth 使用默认 Logger 输出日志。
// func WarnfDepth(depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.outputf("", depth+defaultLoggerDepth, WarnLevel, format, args...)
// 	}
// }

// // WarnTrace 使用默认 Logger 输出日志。
// func WarnTrace(traceID string, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.output(traceID, defaultLoggerDepth, WarnLevel, args...)
// 	}
// }

// // WarnfTrace 使用默认 Logger 输出日志。
// func WarnfTrace(traceID string, format string, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.outputf(traceID, defaultLoggerDepth, WarnLevel, format, args...)
// 	}
// }

// // WarnDepthTrace 使用默认 Logger 输出日志。
// func WarnDepthTrace(traceID string, depth int, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.output(traceID, depth+defaultLoggerDepth, WarnLevel, args...)
// 	}
// }

// // WarnfDepthTrace 使用默认 Logger 输出日志。
// func WarnfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableWarn {
// 		defaultLogger.outputf(traceID, depth+defaultLoggerDepth, WarnLevel, format, args...)
// 	}
// }

// // IsError 返回默认 Logger 是否金鱼 Error 级别。
// func IsError() bool {
// 	return defaultLogger.enableError
// }

// // Error 使用默认 Logger 输出日志。
// func Error(args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.output("", defaultLoggerDepth, ErrorLevel, args...)
// 	}
// }

// // Errorf 使用默认 Logger 输出日志。
// func Errorf(format string, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.outputf("", defaultLoggerDepth, ErrorLevel, format, args...)
// 	}
// }

// // ErrorDepth 使用默认 Logger 输出日志。
// func ErrorDepth(depth int, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.output("", depth+defaultLoggerDepth, ErrorLevel, args...)
// 	}
// }

// // ErrorfDepth 使用默认 Logger 输出日志。
// func ErrorfDepth(depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.outputf("", depth+defaultLoggerDepth, ErrorLevel, format, args...)
// 	}
// }

// // ErrorTrace 使用默认 Logger 输出日志。
// func ErrorTrace(traceID string, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.output(traceID, defaultLoggerDepth, ErrorLevel, args...)
// 	}
// }

// // ErrorfTrace 使用默认 Logger 输出日志。
// func ErrorfTrace(traceID string, format string, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.outputf(traceID, defaultLoggerDepth, ErrorLevel, format, args...)
// 	}
// }

// // ErrorDepthTrace 使用默认 Logger 输出日志。
// func ErrorDepthTrace(traceID string, depth int, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.output(traceID, depth+defaultLoggerDepth, ErrorLevel, args...)
// 	}
// }

// // ErrorfDepthTrace 使用默认 Logger 输出日志。
// func ErrorfDepthTrace(traceID string, depth int, format string, args ...interface{}) {
// 	if defaultLogger.enableError {
// 		defaultLogger.outputf(traceID, depth+defaultLoggerDepth, ErrorLevel, format, args...)
// 	}
// }

// // Recover 使用默认 Logger 输出日志。
// func Recover(recover interface{}) {
// 	defaultLogger.Recover(recover)
// }
