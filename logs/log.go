package logs

import (
	"fmt"
	"log"
	"os"
	_ "reflect"
)

type LogFormat struct {
	Mirror     string `json:"mirror"`
	URI        string `json:"uri"`
	Timestamp  int64  `json:"timestamp"`
	StatusCode int    `json:"status_code"`
	Hit        string `json:"hit"`
	Client     string `json:"client"`
}

type LogFile struct {
	FilePath string
	console  bool
	fd       *os.File
}

type LogKafka struct {
}

type LogRedis struct {
}

type LogAMQP struct {
}

type LogHandler struct {
	File  map[string]string
	Kafka map[string]string
	Redis map[string]string
	AMQP  map[string]string
}

// Info info
func Info(args ...interface{}) {
	log.Println(args)
}

// Debug debug
func Debug(args ...interface{}) {
	log.Println("[Debug]", args)
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Panic(args ...interface{}) {
	log.Panicln(args)
}

var Loggers []Logger

// AccessLog 记录日志
func AccessLog(log LogFormat) {
	for _, logger := range Loggers {
		logger.Write(log)
	}
}

func NewLogger(config map[string]string, name string) Logger {
	var logger Logger
	switch name {
	case "file":
		logger = &LogFile{}
	}
	logger.Connect(config)
	return logger
}

type Logger interface {
	Write(log LogFormat) error
	Connect(config map[string]string)
}

func (this LogFormat) ToString() string {

	return fmt.Sprintf("%d %s %s %s %s %d",
		this.Timestamp,
		this.Client,
		this.Hit,
		this.Mirror,
		this.URI,
		this.StatusCode)
}
