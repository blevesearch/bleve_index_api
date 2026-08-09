//	Copyright (c) 2026 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package log

import "log"

// no logger by default
var globalLogger Logger

// different log levels that can be used by the application layer when logging messages
// note: ordered in a descreasing order of severity
type LogLevel int

const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
)

type Logger interface {
	// level: the verbosity level
	// offset: the position within the calling stack from which the message
	// 		originated. This is useful for contextual loggers which retrieve file/line
	// 		information.
	Log(level LogLevel, offset int, format string, v ...interface{}) error
}

// to be used by the application layer to set a custom logger which implements the Logger interface
func SetLogger(logger Logger) {
	globalLogger = logger
}

func logExf(level LogLevel, offset int, format string, v ...interface{}) {
	if globalLogger != nil {
		err := globalLogger.Log(level, offset+1, format, v...)
		if err != nil {
			log.Printf("Logger error occurred (%s)\n", err)
		}
	}
}

func Debugf(format string, v ...interface{}) {
	logExf(LogDebug, 1, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logExf(LogWarn, 1, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logExf(LogError, 1, format, v...)
}

func Infof(format string, v ...interface{}) {
	logExf(LogInfo, 1, format, v...)
}
