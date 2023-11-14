package logger

import (
	"fmt"
	"github.com/dyammarcano/eventReceiverQuel/internal/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	LogDir      string `yaml:"log_dir"`
	LogToStdout bool   `yaml:"log_to_stdout"`
	AppName     string `yaml:"app_name"`
}

func NewConfig(logDir *string, appName string, stdout bool) *Config {
	if *logDir == "" {
		stdout = true
	}

	return &Config{
		LogDir:      *logDir,
		LogToStdout: stdout,
		AppName:     util.RemoveExtension(appName),
	}
}

func InitLogger(cfg *Config) error {
	var writer zapcore.WriteSyncer
	if cfg.LogToStdout {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	} else {
		if err := util.CreateDirIfNotExists(cfg.LogDir); err != nil {
			log.Printf("cannot create logs directory: %s", cfg.LogDir)
			os.Exit(1)
		}

		outputLogs := filepath.Join(cfg.LogDir, fmt.Sprintf("%s.log", cfg.AppName))

		file, err := os.OpenFile(outputLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = zapcore.AddSync(file)
	}

	productionConfig := zap.NewProductionConfig()
	productionConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(productionConfig.EncoderConfig),
		writer,
		productionConfig.Level,
	)

	zap.ReplaceGlobals(zap.New(core))
	return nil
}

func GetLogger() *zap.Logger {
	return zap.L()
}

func GetSugaredLogger() *zap.SugaredLogger {
	return zap.S()
}

func GetLoggerWithFields(fields ...zap.Field) *zap.Logger {
	return zap.L().With(fields...)
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func Errorf(msg string, fields ...any) {
	zap.S().Errorf(msg, fields...)
}

func Infof(msg string, fields ...any) {
	zap.S().Infof(msg, fields...)
}

func ErrorAndStdout(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
	fmt.Println(msg)
}

func InfoAndStdout(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
	fmt.Println(msg)
}

func InfoAndStdoutf(msg string, fields ...any) {
	zap.S().Infof(msg, fields...)
	fmt.Printf(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}

func Fatalf(msg string, fields ...any) {
	zap.S().Fatalf(msg, fields...)
}
