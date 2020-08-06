package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/spf13/viper"
	"os"
)

var logger *zap.Logger
func init() {
	hook := lumberjack.Logger{
		Filename: viper.GetString("app.logPath"),
		MaxSize: 1,
		MaxBackups: 3,
		MaxAge: 3,
		Compress: false,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey: "time",
		LevelKey: "level",
		NameKey: "logger",
		CallerKey: "caller",
		MessageKey: "msg",
		StacktraceKey: "stacktrace",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook)),atomicLevel,
		)
	caller := zap.AddCaller()
	development := zap.Development()
	field := zap.Fields(zap.String("app","聊天服务器"))
	logger = zap.New(core,caller,development,field)
}

func Logger() *zap.Logger {
	return logger
}

