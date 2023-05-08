package logger

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/beeemT/Packages/fileutil"
)

type Destination int

const (
	File Destination = iota
	Stdout
	Stderr
	Stdin
)

type Logger struct {
	d  Destination
	f  int
	c  int
	q  chan string
	w  *bufio.Writer
	wG *sync.WaitGroup
}

func NewLogger(d Destination, size int, path string, flushAfter int) (*Logger, error) {
	var w *bufio.Writer
	switch d {
	case File:
		path, err := fileutil.PathToAbsFile(path)
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, err
		}
		w = bufio.NewWriter(f)
	case Stdout:
		w = bufio.NewWriter(os.Stdout)
	case Stderr:
		w = bufio.NewWriter(os.Stderr)
	case Stdin:
		w = bufio.NewWriter(os.Stdin)
	}

	q := make(chan string, size)
	wG := &sync.WaitGroup{}
	wG.Add(1)
	l := &Logger{w: w, q: q, d: d, wG: wG, f: flushAfter}
	go l.consumePipe()
	return l, nil
}

func (l *Logger) consumePipe() {
fl:
	for {
		select {
		case v, ok := <-l.q:
			if !ok {
				break fl
			}

			_, err := fmt.Fprintf(l.w, "%s", v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
			l.c++

			if l.c == l.f {
				err = l.w.Flush()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err.Error())
				}
				l.c = 0
			}
		}
	}
	l.wG.Done()
	err := l.w.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
	}
}

func (l *Logger) Logf(format string, a ...interface{}) {
	l.q <- fmt.Sprintf(format, a...)
}

func (l *Logger) Shutdown() {
	close(l.q)
}
