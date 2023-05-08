package logger

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/beeemT/Packages/fileutil"
)

const (
	lipsum = `Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum.`
)

var (
	testFile         = "/home/bt/go/test/logger/log"
	referenceLogFile = "/home/bt/go/test/logger/reference_log"
)

func TestLogger(t *testing.T) {
	f, err := os.Create(testFile)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	err = f.Close()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	var content string
	if flag, _ := fileutil.Exists(referenceLogFile); flag {
		b, err := ioutil.ReadFile(referenceLogFile)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
		content = string(b)
	}

	if content == "" {
		content = lipsum
	}

	l, err := NewLogger(File, 10, testFile, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	l.Logf("%s", content)
	l.Shutdown()
	time.Sleep(time.Second)

	b, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	t.Logf("%s: %s\n", "content", content)
	t.Logf("%s: %s\n", "written", string(b))

	if string(b) != content {
		t.Log("Failed because the logged content did not match the reference content.")
		t.Fail()
	}
}
