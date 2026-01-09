package main

import (
	"fmt"
	"time"

	logger "example.com/lesson071225/logger"
)

type IReader interface {
	Read() string
}

type IWriter interface {
	Write(s string)
}

type IReaderWriter interface {
	IReader
	IWriter
}

type MemoryBuffer struct {
	data string
}

func (mb *MemoryBuffer) Read() string {
	return mb.data
}

func (mb *MemoryBuffer) Write(s string) {
	mb.data += s + "\n"
}

type LogMessage struct {
	UserID int
	msg    string
}

func (lm LogMessage) String() string {
	lm.msg = fmt.Sprintf("[%s] UserID: %d msg: %s", time.Now().String(), lm.UserID, lm.msg)
	return lm.msg
}

func main() {

	var l logger.ILogger
	l.LogInfo(LogMessage{UserID: 123, msg: "This is an info message"})
	l.LogWarning(LogMessage{UserID: 2, msg: "This is a warning message"})
	l.LogError(LogMessage{UserID: 3, msg: "This is an error message"})

	var b IReaderWriter = &MemoryBuffer{}
	b.Write("Hello World!")
	b.Write("asdasd")
	b.Write("Go to GOlang!")
	fmt.Println(b.Read())
}
