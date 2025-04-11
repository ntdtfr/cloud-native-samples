// pkg/logger/logger.go
package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields map[string]interface{}

type Logger interface {
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, err error, fields ...Fields)
	Error(msg string, err error, fields ...Fields)
	Fatal(msg string, err error, fields ...Fields)
}

type zapLogger struct {
	logger *zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(logLevel string) Logger {
	// Parse log level
	level := zap.InfoLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore
    EncodeLevel:    zapcore.LowercaseLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
  }

  // Create encoder
  encoder := zapcore.NewJSONEncoder(encoderConfig)

  // Create core
  core := zapcore.NewCore(
    encoder,
    zapcore.AddSync(os.Stdout),
    level,
  )

  // Create logger
  logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

  return &zapLogger{
    logger: logger,
  }
}

func (l *zapLogger) Debug(msg string, fields ...Fields) {
  l.logger.Debug(msg, l.getZapFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Fields) {
  l.logger.Info(msg, l.getZapFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, err error, fields ...Fields) {
  fieldsWithError := l.appendError(err, fields...)
  l.logger.Warn(msg, l.getZapFields(fieldsWithError...)...)
}

func (l *zapLogger) Error(msg string, err error, fields ...Fields) {
  fieldsWithError := l.appendError(err, fields...)
  l.logger.Error(msg, l.getZapFields(fieldsWithError...)...)
}

func (l *zapLogger) Fatal(msg string, err error, fields ...Fields) {
  fieldsWithError := l.appendError(err, fields...)
  l.logger.Fatal(msg, l.getZapFields(fieldsWithError...)...)
}

// getZapFields converts a variadic list of Fields into a slice of zap.Field.
// It iterates over each Field map, and for each key-value pair, it creates
// a zap.Field using zap.Any. The resulting slice of zap.Field is returned.
func (l *zapLogger) getZapFields(fields ...Fields) []zap.Field {
  var zapFields []zap.Field

  for _, field := range fields {
    for k, v := range field {
      zapFields = append(zapFields, zap.Any(k, v))
    }
  }

  return zapFields
}

// appendError appends an error to a given list of Fields. If the given list of Fields is empty, it returns a new list with the error as a single field.
// If the given list of Fields is not empty, it appends the error to the first field.
func (l *zapLogger) appendError(err error, fields ...Fields) []Fields {
  if err == nil {
    return fields
  }

  if len(fields) == 0 {
    return []Fields{{"error": err.Error()}}
  }

  fields[0]["error"] = err.Error()
  return fields
}
