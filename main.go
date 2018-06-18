package main

import (
	"flag"
	"github.com/beengo/mirror-cache/logs"
)

// InitLoggers 初始化日志
func initLoggers() {

	if len(Config.Logger.File) > 0 {
		logger := logs.NewLogger(Config.Logger.File, "file")
		logs.Loggers = append(logs.Loggers, logger)
	}
	// if len(Config.Logger.Redis) > 0 {
	// 	logs.NewLogger(Config.Logger.Redis, "redis")
	// }

}

func main() {
	configFile := flag.String("cfg", "config.yml", "config file")
	flag.Parse()
	LoadConfig(*configFile)
	initLoggers()
	logs.Info("Starting Server...")
	StartHttpServer()
}
