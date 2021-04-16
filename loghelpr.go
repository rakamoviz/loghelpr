package loghelpr

import (
	"context"

	"github.com/rakamoviz/loghelpr/contextkeys"
	"github.com/sirupsen/logrus"
)

var (
	std = logrus.NewEntry(logrus.StandardLogger())
)

const (
	// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
	// ensure the formatted time is always the same number of characters.
	RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"
)

func BuildContext(ctx context.Context, logger *logrus.Entry, sampleLog bool) context.Context {
	return context.WithValue(
		context.WithValue(ctx, contextkeys.Logger{}, logger),
		contextkeys.SampleLog{}, sampleLog,
	)
}

func Get(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(contextkeys.Logger{})

	if logger == nil {
		return std
	}

	return logger.(*logrus.Entry)
}

func Fn(ctx context.Context, functionName string, args *map[string]interface{}, forceLogs ...bool) func() {
	forceLog := forceLog(forceLogs...)

	logger := Get(ctx).WithFields(logrus.Fields{
		"fnEntrance": 1,
		"fnName":     functionName,
		"fnArgs":     args,
	})

	log(ctx, logger, forceLog, logrus.InfoLevel, nil)

	return func() {
		logger := Get(ctx).WithFields(logrus.Fields{
			"fnExit": 1,
			"fnName": functionName,
		})

		log(ctx, logger, forceLog, logrus.InfoLevel, nil)
	}
}

func log(ctx context.Context, logger *logrus.Entry, forceLog bool, level logrus.Level, args ...interface{}) {
	sampleLog, _ := ctx.Value(contextkeys.SampleLog{}).(bool)

	if sampleLog || forceLog {
		logger.Log(level, args...)
	}
}

func forceLog(forceLogs ...bool) bool {
	if len(forceLogs) == 0 {
		return false
	}

	return forceLogs[0]
}

func Trace(ctx context.Context, args ...interface{}) func(...bool) {
	return func(forceLogs ...bool) {
		log(ctx, Get(ctx), forceLog(forceLogs...), logrus.TraceLevel, args...)
	}
}

func Debug(ctx context.Context, args ...interface{}) func(...bool) {
	return func(forceLogs ...bool) {
		log(ctx, Get(ctx), forceLog(forceLogs...), logrus.DebugLevel, args...)
	}
}

func Info(ctx context.Context, args ...interface{}) func(...bool) {
	return func(forceLogs ...bool) {
		log(ctx, Get(ctx), forceLog(forceLogs...), logrus.InfoLevel, args...)
	}
}

func Warn(ctx context.Context, args ...interface{}) func(...bool) {
	return func(forceLogs ...bool) {
		log(ctx, Get(ctx), forceLog(forceLogs...), logrus.WarnLevel, args...)
	}
}

func Warning(ctx context.Context, args ...interface{}) func(...bool) {
	return Warn(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) func() {
	return func() {
		log(ctx, Get(ctx), true, logrus.ErrorLevel, args...)
	}
}

func Fatal(ctx context.Context, args ...interface{}) func() {
	return func() {
		log(ctx, Get(ctx), true, logrus.FatalLevel, args...)
	}
}

func Panic(ctx context.Context, args ...interface{}) func() {
	return func() {
		log(ctx, Get(ctx), true, logrus.PanicLevel, args...)
	}
}
