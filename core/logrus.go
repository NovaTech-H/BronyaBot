package core

import (
	"BronyaBot/global"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

// Format formats the log entry to include colored level indicators and caller information
func (t LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	log := global.Config.Logger
	timeStamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		_, _ = fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", log.Prefix, timeStamp, levelColor, entry.Level, funcVal, fileVal, entry.Message)
	} else {
		_, _ = fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s\n", log.Prefix, timeStamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	mLog := logrus.New()
	mLog.SetOutput(os.Stdout)
	mLog.SetReportCaller(global.Config.Logger.ShowLine)
	mLog.SetFormatter(&LogFormatter{})
	level, err := logrus.ParseLevel(global.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	mLog.SetLevel(level)
	scheduleLogRotation(mLog)
	// 第一次运行时创建日志文件
	fileWriter := createLogFile()
	if fileWriter != nil {
		mLog.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
	}
	return mLog
}

func scheduleLogRotation(logger *logrus.Logger) {
	go func() {
		// 计算到下一个 00:00 的时间
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)

		// 创建一个定时器，在下一个 00:00 触发
		timer := time.NewTimer(durationUntilMidnight)
		<-timer.C

		// 在 00:00 创建新的日志文件
		fileWriter := createLogFile()
		if fileWriter != nil {
			logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
		}

		// 之后每隔 24 小时触发一次
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			fileWriter := createLogFile()
			if fileWriter != nil {
				logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
			}
		}
	}()
}

//func scheduleLogRotation(logger *logrus.Logger) {
//	go func() {
//		// 计算当前时间到下一个 00:00 的时间间隔
//		now := time.Now()
//		next := now.Add(24 * time.Hour)
//		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
//		duration := next.Sub(now)
//
//		// 等待到下一个 00:00
//		time.Sleep(duration)
//
//		// 创建日志文件
//		fileWriter := createLogFile()
//		if fileWriter != nil {
//			logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
//		}
//
//		// 每隔 24 小时触发一次
//		for range time.NewTicker(24 * time.Hour).C {
//			fileWriter := createLogFile()
//			if fileWriter != nil {
//				logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
//			}
//		}
//	}()
//}

func createLogFile() *os.File {
	logDir := global.Config.Logger.Director
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return nil
	}

	//var logFileName = fmt.Sprintf("%s/%s.log", logDir, time.Now().Format("2006-01-02_15-04-05"))
	var logFileName = fmt.Sprintf("%s/%s.log", logDir, time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return nil
	}
	return logFile
}
