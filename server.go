package main

import (
	"github.com/beengo/mirror-cache/logs"
	"net/http"
	"strings"
)

var staticHandlers map[string]http.Handler

func accessLog(req *http.Request) {
	logs.Info(req.RemoteAddr, req.RequestURI)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	go accessLog(req)
	mirrorFound := false
	for _, mirror := range Config.Mirrors {
		if strings.HasPrefix(req.RequestURI, mirror.Prefix) {
			mirrorFound = true
			CacheBy(mirror, w, req)
		}
	}
	if !mirrorFound {
		w.WriteHeader(404)
	}
}

func initStaticHandlers() {
	staticHandlers = make(map[string]http.Handler)
	for _, mirror := range Config.Mirrors {
		logs.Debug("init static handler for", mirror.Name)
		staticHandlers[mirror.Name] = http.FileServer(http.Dir(mirror.LocalDir))
	}
}

func StartHttpServer() {
	http.HandleFunc("/", handleAll)
	initStaticHandlers()
	http.ListenAndServe(Config.Listen, nil)
	logs.Info("Listen at", Config.Listen)
	logs.Info("Block size", Config.BlockSize)
	select {}
}
