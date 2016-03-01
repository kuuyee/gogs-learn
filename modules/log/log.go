package log

import (
	"fmt"
	"sync"
)

var (
	loggers []*Logger
)

func Error(skip int, format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Close() //方法调用错误，为了调试临时使用
	}
}

func Fatal(skip int, format string, v ...interface{}) {
	Error(skip, format, v...)
	for _, l := range loggers {
		l.Close()
	}
}

type LogerInterface interface {
	Init(config string) error
	WriteMsg(msg string, skip, level int) error
	Destroy()
	Flush()
}

type logMsg struct {
	skip, level int
	msg         string
}

type Logger struct {
	adapter string
	lock    sync.Mutex
	level   int
	msg     chan *logMsg
	outputs map[string]LogerInterface
	quit    chan bool
}

// Close 关闭logger,刷新所有的通道数据和销毁所有adapter实例
func (l *Logger) Close() {
	l.quit <- true
	for {
		if len(l.msg) > 0 {
			bm := <-l.msg
			for _, l := range l.outputs {
				if err := l.WriteMsg(bm.msg, bm.skip, bm.level); err != nil {
					fmt.Println("ERROR, unable to WriteMsg:", err)
				}
			}
		} else {
			break
		}
	}
}
