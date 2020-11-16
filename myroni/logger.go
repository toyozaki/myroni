package myroni

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

const defaultLogFormat = "{{.StartTime}} {{.ElapsedTime}} {{.Src}} {{.Path}}"

type MyLog struct {
	StartTime   time.Time
	ElapsedTime time.Duration
	Src         string
	Path        string
}

type MyLogger interface {
	Print(v ...interface{})
	Println(v ...interface{})
}

type Logger struct {
	MyLogger
	dateFormat string
	template   *template.Template
}

func NewLogger() *Logger {
	return &Logger{
		MyLogger:   log.New(os.Stdout, "[myroni] ", 0),
		dateFormat: time.RFC3339,
		template:   template.Must(template.New("myroni_parser").Parse(defaultLogFormat)),
	}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	myLog := MyLog{
		StartTime:   start,
		ElapsedTime: time.Since(start),
		Src:         r.RemoteAddr,
		Path:        r.URL.Path,
	}

	buff := &bytes.Buffer{}
	l.template.Execute(buff, myLog)
	l.Println(buff)
}
