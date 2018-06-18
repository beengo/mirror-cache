package main

import (
	"github.com/noaway/dateparse"
	"strings"
	"testing"
)

func TestLocalFileExists(t *testing.T) {
	uri := "/ab/c?ui=xx"
	expect := "data/test/ab/c"
	path, err := localFilePath("data/test", uri)
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(path, expect) != 0 {
		t.Error("Expect: ", expect, "got: ", path)
	}
}

func TestLastModifiedTime(t *testing.T) {
	timeStr := "Wed, 11 Oct 2006 09:10:50 GMT"
	_, err := dateparse.ParseAny(timeStr)
	if err != nil {
		t.Error(err)
	}
}
