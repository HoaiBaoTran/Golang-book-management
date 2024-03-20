package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	FileLogger    *zap.SugaredLogger
	ConsoleLogger *zap.SugaredLogger
}

func InitLogger() Logger {
	fileLogger := InitFileLogger()
	consoleLogger := InitConsoleLogger()
	return Logger{
		ConsoleLogger: consoleLogger,
		FileLogger:    fileLogger,
	}
}

func InitFileLogger() *zap.SugaredLogger {
	writeSync := getLogWriter()
	encoder := getEncoder()

	zapCore := zapcore.NewCore(encoder, writeSync, zapcore.DebugLevel)
	fileLogger := zap.New(zapCore, zap.AddCaller())
	return fileLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	/*
		* NewJSONEncoder
		{
			"level":"[INFO]",
			"time":"2021-07-17 14:54:05",
			"caller":"golangforteaching/main.go:10",
			"message":"Today is :2021-July-17"
		}

		* NewConsoleEncoder
		* 2021-07-17 14:52:35	[INFO]	golangforteaching/main.go:10	Today is :2021-July-17
	*/
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		TimeKey:      "time",
		CallerKey:    "caller",
		LevelKey:     "level",
		EncodeLevel:  customLevelEncoder,
		EncodeTime:   syslogTimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})
}

func getLogWriter() zapcore.WriteSyncer {
	current_date := time.Now().Format("02-01-2006")
	current_time := time.Now().Format("15-04-05")
	outputFileName := fmt.Sprintf("pkg/logger/logger-files/log_%s_%s.txt", current_date, current_time)
	file, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Can't open log file", err)
		// log.Fatal("Can't open log file", err)
	}
	return zapcore.AddSync(file)
}

func InitConsoleLogger() *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:    "console", // json or console
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder,
			EncodeLevel:  customLevelEncoder,
			EncodeTime:   syslogTimeEncoder,
		},
	}
	consoleLogger, err := cfg.Build()
	if err != nil {
		log.Fatal("Can't not create logger", err)
	}
	return consoleLogger.Sugar()
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
