package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init(){
	initLogger()
}

func initLogger() {
    cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString(t.Format("2006-01-02 15:04:05.999999999")) 
    }
	cfg.DisableCaller = true
    cfg.Encoding = "console"
    cfg.OutputPaths = []string{"stdout","chat.log"}
    l, _ := cfg.Build()
    logger = l
}

func Logi(a ...any) {
	logger.Sugar().Info(a)
}

func Loge(a ...any) {
	logger.Sugar().Error(a)
}

