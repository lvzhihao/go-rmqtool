package rmqtool

import "log"

var (
	Log *Logger
)

type Logger struct {
	instance interface{}
}

type LoggerInterface interface {
	Error(...interface{})
	Debug(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
	Warn(...interface{})
	Info(...interface{})
}

func init() {
	Log = &Logger{
		instance: &DefaultLog{},
	}
}

func (c *Logger) Error(v ...interface{}) {
	c.instance.(LoggerInterface).Error(v...)
}

func (c *Logger) Debug(v ...interface{}) {
	c.instance.(LoggerInterface).Debug(v...)
}

func (c *Logger) Info(v ...interface{}) {
	c.instance.(LoggerInterface).Info(v...)
}

func (c *Logger) Warn(v ...interface{}) {
	c.instance.(LoggerInterface).Warn(v...)
}

func (c *Logger) Panic(v ...interface{}) {
	c.instance.(LoggerInterface).Panic(v...)
}

func (c *Logger) Fatal(v ...interface{}) {
	c.instance.(LoggerInterface).Fatal(v...)
}

type DefaultLog struct {
	LoggerInterface
}

func (c *DefaultLog) Error(v ...interface{}) {
	log.Println(append([]interface{}{"Error"}, v...))
}

func (c *DefaultLog) Debug(v ...interface{}) {
	log.Println(append([]interface{}{"Debug"}, v...))
}

func (c *DefaultLog) Warn(v ...interface{}) {
	log.Println(append([]interface{}{"Warn"}, v...))
}

func (c *DefaultLog) Info(v ...interface{}) {
	log.Println(append([]interface{}{"Info"}, v...))
}

func (c *DefaultLog) Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func (c *DefaultLog) Panic(v ...interface{}) {
	log.Panic(v...)
}
