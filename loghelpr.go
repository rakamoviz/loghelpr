/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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

	// TextFormat represents the text logging format
	TextFormat = "text"

	// JSONFormat represents the JSON logging format
	JSONFormat = "json"
)

// WithLogger returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func BuildContext(ctx context.Context, logger *logrus.Entry, sampleLog bool) context.Context {
	return context.WithValue(
		context.WithValue(ctx, contextkeys.Logger{}, logger),
		contextkeys.SampleLog{}, sampleLog,
	)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func Get(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(contextkeys.Logger{})

	if logger == nil {
		return std
	}

	return logger.(*logrus.Entry)
}

func LogFn(ctx context.Context, functionName string, args *map[string]interface{}) func(...bool) string {
	return func(forceLogs ...bool) string {
		logger := Get(ctx).WithFields(logrus.Fields{
			"fnExit": 1,
			"fnName": functionName,
		})

		log(ctx, logger, forceLog(forceLogs...), logrus.InfoLevel, nil)

		return functionName
	}
}

func Fn(ctx context.Context, functionName string, args *map[string]interface{}) func(...bool) (context.Context, string, *map[string]interface{}) {
	return func(forceLogs ...bool) (context.Context, string, *map[string]interface{}) {
		logger := Get(ctx).WithFields(logrus.Fields{
			"fnEntrance": 1,
			"fnName":     functionName,
			"fnArgs":     args,
		})

		log(ctx, logger, forceLog(forceLogs...), logrus.InfoLevel, nil)

		return ctx, functionName, args
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
