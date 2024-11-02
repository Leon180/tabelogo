package utility

import (
	"context"
	"errors"
	"os"
	"time"

	"authenticate/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
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

type GormLogger struct {
	zap *zap.Logger
	gormlogger.Config
}

func NewGormZapLogger(zap *zap.Logger) *GormLogger {
	return &GormLogger{
		zap:    zap,
		Config: gormlogger.Config{IgnoreRecordNotFoundError: true},
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.zap.Sugar().Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.zap.Sugar().Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.zap.Sugar().Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		l.zap.Error("trace",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn {
		l.zap.Warn("trace",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.String("slow-threshold", l.SlowThreshold.String()),
		)
		return
	}

	if l.LogLevel >= gormlogger.Info {
		l.zap.Info("trace",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	}
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
