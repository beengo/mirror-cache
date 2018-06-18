package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (this *LogFile) Connect(config map[string]string) {
	this.FilePath = config["filepath"]
	this.console = config["console"] == "true"
	err := os.MkdirAll(filepath.Dir(this.FilePath), 0666)
	if err != nil {
		Fatal("Failed to create dir", filepath.Dir(this.FilePath))
	}
	this.fd, err = os.OpenFile(this.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Panic("Failed to open log file: ", this.FilePath, err)
	}
	Debug("Ready to log to file", this.FilePath)
}

func (this *LogFile) Write(log LogFormat) error {
	logs, err := json.Marshal(log)
	if err != nil {
		return err
	}
	logStr := string(logs)
	_, err = this.fd.WriteString(logStr)
	if this.console {
		fmt.Println(logStr)
	}
	return err
}
