package log

import (
	"Share-Wallet/pkg/common/config"
	"bufio"
	"fmt"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	logger2 "gorm.io/gorm/logger"
)

var logger *Logger
var sqlLogger *Logger

type Logger struct {
	*logrus.Logger
	Pid int
}

type Writer struct {
}

func init() {
	logger = loggerInit("")

}
func NewPrivateLog(moduleName string) {
	logger = loggerInit(moduleName)
}

func PrivateLog(moduleName string, logLevel uint32, rotationTime int, rotationCount uint, elasticSearchSwitch bool) {
	logger = newloggerInit(moduleName, logLevel, rotationTime, rotationCount, elasticSearchSwitch)
}

//func GetNewLogger(moduleName string) *Logger {
//	sqlLogger = loggerInit(moduleName)
//	return sqlLogger
//}

func GetSqlLogger(moduleName string) logger2.Interface {
	//0 panic 1 fetal 2 error 3 warn 4 info 5 debug 6 trace
	sqlLogger = loggerInit(moduleName)
	var sqlLogLevel logger2.LogLevel
	switch config.Config.Log.GormLogLevel {
	case 1:
		sqlLogLevel = logger2.Silent
	case 2:
		sqlLogLevel = logger2.Error
	case 3:
		sqlLogLevel = logger2.Warn
	case 4:
		sqlLogLevel = logger2.Info
	default:
		sqlLogLevel = logger2.Error
	}

	newLogger := logger2.New(
		Writer{},
		logger2.Config{
			SlowThreshold:             3000 * time.Millisecond,
			LogLevel:                  sqlLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	return newLogger
}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	sqlLogger.Infof(format, args...)
}

func loggerInit(moduleName string) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.Level(config.Config.Log.RemainLogLevel))

	//Close std console output
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(src)
	logger.SetOutput(writer)
	//logger.SetOutput(os.Stdout)
	//Log Console Print Style Setting
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	//File name and line number display hook
	logger.AddHook(newFileHook())

	//Send logs to elasticsearch hook
	if config.Config.Log.ElasticSearchSwitch {
		graylogHook := graylog.NewGraylogHook(config.Config.Log.ElasticSearchAddr[0], map[string]interface{}{"this": moduleName})
		logger.AddHook(graylogHook)
	} else {
		//Log file segmentation hook
		hook := NewLfsHook(time.Duration(config.Config.Log.RotationTime)*time.Hour, config.Config.Log.RemainRotationCount, moduleName)
		logger.AddHook(hook)
	}

	return &Logger{
		logger,
		os.Getpid(),
	}
}

func newloggerInit(moduleName string, logLevel uint32, rotationTime int, rotationCount uint, elasticSearchSwitch bool) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.Level(logLevel))

	//Close std console output
	if moduleName != "" {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err.Error())
		}
		writer := bufio.NewWriter(src)
		logger.SetOutput(writer)
	}

	//logger.SetOutput(os.Stdout)
	//Log Console Print Style Setting
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	//File name and line number display hook
	logger.AddHook(newFileHook())

	//Send logs to elasticsearch hook
	if elasticSearchSwitch {
		graylogHook := graylog.NewGraylogHook(config.Config.Log.ElasticSearchAddr[0], map[string]interface{}{"this": moduleName})
		logger.AddHook(graylogHook)
	} else {
		//Log file segmentation hook
		if moduleName != "" {
			hook := NewLfsHook(time.Duration(rotationTime)*time.Hour, rotationCount, moduleName)
			logger.AddHook(hook)
		}
	}

	return &Logger{
		logger,
		os.Getpid(),
	}
}

func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	return lfsHook
}
func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, moduleName string) *rotatelogs.RotateLogs {
	if moduleName != "" {
		moduleName = moduleName + "."
	}

	location := "../logs/"
	if config.Config.Log.StorageLocation != "" {
		location = config.Config.Log.StorageLocation
	}
	writer, err := rotatelogs.New(
		location+moduleName+level+"."+"%Y-%m-%d",
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

func Info(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}

func Error(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}

func Debug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}

func Warn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}

// Deprecated
func Warning(token, OperationID, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Warningf(format, args...)

}

// Deprecated
func InfoByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args)
}

// Deprecated
func ErrorByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}

// Print log information in k, v format,
// kv is best to appear in pairs. tipInfo is the log prompt information for printing,
// and kv is the key and value for printing.
// Deprecated
func InfoByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Info(tipInfo)
}

// Deprecated
func ErrorByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Error(tipInfo)
}

// Deprecated
func DebugByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Debug(tipInfo)
}

// Deprecated
func WarnByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Warn(tipInfo)
}

// internal method
func argsHandle(OperationID string, fields logrus.Fields, args []interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
	fields["OperationID"] = OperationID
	fields["PID"] = logger.Pid
}
func NewInfo(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}
func NewError(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}
func NewDebug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}
func NewWarn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}

//func (l *Logger) LogMode(level LogLevel) logger2.Interface {
//	newlogger := *l
//	logger2.LogLevel = level
//	var logrusLevel logrus.Level
//	switch level {
//	case logger2.Silent:
//		logrusLevel = logrus.PanicLevel
//	case logger2.Info:
//		logrusLevel = logrus.InfoLevel
//	case logger2.Warn:
//		logrusLevel = logrus.WarnLevel
//	case logger2.Error:
//		logrusLevel = logrus.ErrorLevel
//	default:
//		logrusLevel = logrus.ErrorLevel
//	}
//	sqlLogger.SetLevel(logrusLevel)
//	return sqlLogger
//}
//
//func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
//	sqlLogger.WithFields(logrus.Fields{
//		"OperationID": "",
//		"PID":         sqlLogger.Pid,
//	}).Infoln(s, i)
//}
//
//func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
//	sqlLogger.WithFields(logrus.Fields{
//		"OperationID": "",
//		"PID":         sqlLogger.Pid,
//	}).Warnln(s, i)
//}
//
//func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
//	sqlLogger.WithFields(logrus.Fields{
//		"OperationID": "",
//		"PID":         sqlLogger.Pid,
//	}).Errorln(s, i)
//}
//
//func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
//	var errorString string
//	if err != nil {
//		errorString = err.Error()
//	}
//
//	sql, rowsAffected := fc()
//	sqlLogger.WithFields(logrus.Fields{
//		"OperationID": "",
//		"PID":         sqlLogger.Pid,
//	}).Infoln("sql:", sql, "rowsAffected:", rowsAffected, "error:", errorString)
//}
