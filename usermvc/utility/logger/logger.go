package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"usermvc/globals"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type key string

const (
	_callerKey = "caller"
	_requestID = "requestID"
)

// GetLoggerWithContext returns a global logger with Proper CallerName and TranceID
func GetLoggerWithContext(ctx context.Context) *zap.SugaredLogger {
	if getrequestID(ctx) == "" {
		ctx = SetRequestID(ctx)
	}
	return zap.S().With(_requestID, getrequestID(ctx))
}

func getrequestID(ctx context.Context) string {
	traceID := ctx.Value(key(_requestID))
	if traceID != nil {
		return traceID.(string)
	}
	fmt.Println("traceID is", traceID)
	return ""
}

func SetRequestID(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, key(_requestID), generateReuestID())
	return ctx
}

func generateReuestID() string {
	traceID, _ := uuid.NewV4()
	fmt.Println("printing requestid", traceID)
	return traceID.String()

}

func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()
	conf := globals.GetConfig()
	logFile := filepath.Join("", conf.FileName)
	if _, err := os.Create(logFile); err != nil {
		logger.Panic("error while creating log files")
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	})
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config.EncoderConfig), w, zap.InfoLevel)
	logger = zap.New(core, zap.AddCaller())
	logger.Info("starting logger")
	return logger
}
