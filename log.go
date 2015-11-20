// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mssfcore

import (
	"strings"

	"github.com/cincout/mssfcore/logs"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	MssfLogger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	MssfLogger.EnableFuncCallDepth(b)
	MssfLogger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
var MssfLogger *logs.MssfLogger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := MssfLogger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(v ...interface{}) {
	MssfLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	MssfLogger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	MssfLogger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	MssfLogger.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	MssfLogger.Warning(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Warn(v ...interface{}) {
	MssfLogger.Warn(generateFmtStr(len(v)), v...)
}

func Notice(v ...interface{}) {
	MssfLogger.Notice(generateFmtStr(len(v)), v...)
}

// Info logs a message at info level.
func Informational(v ...interface{}) {
	MssfLogger.Informational(generateFmtStr(len(v)), v...)
}

// compatibility alias for Warning()
func Info(v ...interface{}) {
	MssfLogger.Info(generateFmtStr(len(v)), v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	MssfLogger.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	MssfLogger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
