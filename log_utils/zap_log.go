package zlogger

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log        *zap.Logger
	logpath    = filepath.Join(os.TempDir(), "tinyurl", fmt.Sprintf("test-log-%d", time.Now().UnixNano()))
	zaplogPath = filepath.Join("winfile:///", logpath)
)

func newWinfileSink(u *url.URL) (zap.Sink, error) {
	var name string
	if u.Path != "" {
		name = u.Path[1:]
	} else if u.Opaque != "" {
		name = u.Opaque[1:]
	} else {
		return nil, errors.New("path error")
	}
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func init() {
	fmt.Println(os.TempDir())
	zap.RegisterSink("winfile", newWinfileSink)
	logConfig := zap.Config{

		OutputPaths: []string{zaplogPath},
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	Log, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
}

func Debug(msg string, tags ...zap.Field) {
	Log.Debug(msg, tags...)
	Log.Sync()
}

func Info(msg string, tags ...zap.Field) {
	Log.Info(msg, tags...)
	Log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	Log.Error(msg, tags...)
	Log.Sync()
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
