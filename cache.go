package main

import (
	"github.com/beengo/mirror-cache/logs"
	"github.com/noaway/dateparse"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CacheBy(mirror Mirror, w http.ResponseWriter, req *http.Request) {
	logs.Debug("Processing by", mirror.Name, ", url: ", req.RequestURI)
	log := logs.LogFormat{
		StatusCode: 200,
		URI:        req.RequestURI,
		Timestamp:  time.Now().Unix(),
		Hit:        "cache",
		Mirror:     mirror.Name,
		Client:     req.RemoteAddr,
	}
	if needDownload(mirror, req.RequestURI) {
		logs.Debug("Need download", req.RequestURI)
		DownloadAndServe(mirror, w, req)
		log.Hit = "remote"
	} else {
		logs.Debug("Serve from cache", req.RequestURI)
		staticHandlers[mirror.Name].ServeHTTP(w, req)
	}

	go logs.AccessLog(log)
	// w.Write([]byte(mirror.Name))
}

func copyHeaders(w http.ResponseWriter, resp *http.Response) {
	var headers = []string{
		"content-type",
		"last-modified",
		"content-length",
	}
	for _, key := range headers {
		value := resp.Header.Get(key)
		if len(value) > 0 {
			w.Header().Set(key, value)
		}
	}
}

func DownloadAndServe(mirror Mirror, w http.ResponseWriter, req *http.Request) {
	localPath, _ := localFilePath(mirror.LocalDir, req.RequestURI)
	url := remoteUrl(mirror, req.RequestURI)
	resp, err := http.Get(url)
	if err != nil {
		logs.Debug(err)
	}
	if resp.StatusCode != 200 {
		logs.Debug("Remote url response not ok", resp.StatusCode, url)
		w.WriteHeader(resp.StatusCode)
		return
	}
	copyHeaders(w, resp)
	defer resp.Body.Close()
	blockSize := Config.BlockSize

	err = os.MkdirAll(filepath.Dir(localPath), 0666)
	if err != nil {
		logs.Fatal("Failed to create dir", filepath.Dir(localPath))
	}
	fp, ferr := os.Create(localPath)
	if ferr != nil {
		logs.Fatal("Failed to open local file", localPath, err)
	}
	var bb []byte = make([]byte, blockSize)
	var readCount int64 = 0
	write := func(bf []byte) {
		w.Write(bf)
		fp.Write(bf)
	}
	for {
		n, err := io.ReadFull(resp.Body, bb)
		if err == io.EOF {
			break
		}
		if err != nil {
			// 读到底了, 正常的, 重置byte
			if resp.ContentLength-readCount == int64(n) {
				write(bb[:n])
			} else { // 读的过程中真有问题了
				log.Fatal("IO Read error", err)
			}
		} else { // 正常的， 继续读
			write(bb)
		}
		readCount += int64(n)
	}
}

func localFilePath(localDir string, reqURI string) (string, error) {
	uri, err := url.ParseRequestURI(reqURI)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(localDir, "/") + uri.Path, nil
}

func remoteUrl(mirror Mirror, reqURI string) string {
	return strings.TrimRight(mirror.ProxyTarget, "/") + reqURI
}

func needDownload(mirror Mirror, reqURI string) bool {
	path, _ := localFilePath(mirror.LocalDir, reqURI)
	finfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		logs.Debug("cached File not found", path)
		return true
	}
	if err != nil {
		logs.Fatal(err)
	}
	now := time.Now().Unix()
	modTime := finfo.ModTime().Unix()
	if modTime+mirror.Expire > now {
		return false
	}
	// 获取服务器的最后更改时间， 拿不到时不下载
	resp, _ := http.Head(remoteUrl(mirror, reqURI))
	lastMod, err := dateparse.ParseAny(resp.Header.Get("last-modified"))
	if err != nil {
		logs.Fatal(err)
		return false
	}
	return modTime < lastMod.Unix()
}
