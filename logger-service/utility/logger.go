package utility

import (
	"context"
	"logger-service/config"
	"logger-service/model/enum"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
	logRecord   *LogRecord
)

func initLogRecord(logRecordDay int, logFileName string, lumberJackLogger *lumberjack.Logger) {
	if logRecord != nil {
		return
	}
	logRecord = &LogRecord{
		logRecordDay:     logRecordDay,
		logFileName:      logFileName,
		lumberJackLogger: lumberJackLogger,
	}
}

type LogRecord struct {
	logRecordDay     int
	logFileName      string
	lumberJackLogger *lumberjack.Logger
}

func (l *LogRecord) RotateIfNeed() {
	currentDay := time.Now().Day()
	if l.logRecordDay == 0 {
		l.logRecordDay = currentDay
		return
	}
	if l.logRecordDay == currentDay {
		return
	}
	l.logRecordDay = currentDay
	l.lumberJackLogger.Rotate()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func InitLogger(logCfg config.LogConfig) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logCfg.FileName,
		MaxSize:    logCfg.MaxSize,
		MaxAge:     logCfg.MaxAge,
		MaxBackups: logCfg.MaxBackups,
		Compress:   logCfg.Compress,
	}
	cfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "severity",
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	initLogRecord(0, logCfg.FileName, lumberJackLogger)
	Logger = zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg),
				zapcore.NewMultiWriteSyncer(
					zapcore.AddSync(lumberJackLogger),
					zapcore.AddSync(os.Stderr),
				),
				logCfg.Level,
			),
			zapcore.NewCore(
				getEncoder(),
				zapcore.AddSync(os.Stdout),
				logCfg.Level,
			),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Hooks(func(e zapcore.Entry) error {
			logRecord.RotateIfNeed()
			return nil
		}),
	)
	SugarLogger = Logger.Sugar()
}

func LogWithTraceID(ctx context.Context, data ...interface{}) {
	if traceID := ctx.Value(enum.TraceIDKey); traceID != nil {
		SugarLogger.
			With("EventID", traceID.(string)).
			Error(data...)
		return
	}
	SugarLogger.Error(data...)
}
