package log

import (
	"bbexgo/help"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	infoLog   *log.Logger
	warnLog   *log.Logger
	errLog    *log.Logger
	fatalLog  *log.Logger
	logInfo   = log.Ltime | log.Ldate | log.Lshortfile
	fileMap   = make(map[string]*os.File)
	logDir    = "logs/"
	debuger   = false
	levelList = []string{"Info", "Warning", "Error", "Fatal"}
	initDaily = time.Now().Format("20060102")
)

func SetDebug(s bool) {
	debuger = s
}

func Info(args ...interface{}) {
	write(getLog("Info"), args...)
}

func Warning(args ...interface{}) {
	write(getLog("Warning"), args...)
}

func Error(args ...interface{}) {
	write(getLog("Error"), args...)
}

func Fatal(args ...interface{}) {
	write(getLog("Fatal"), args...)
	panic(fmt.Sprintln(args...))
}

func getFile(name string) *os.File {
	nowDaily := time.Now().Format("20060102")
	fileName := name + "_" + nowDaily + ".log"
	if f, ok := fileMap[fileName]; ok {
		return f
	}
	file, err := os.OpenFile(getDir()+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file ERROR: ", err)
	}
	return file
}

func getDir() string {
	currPath, err := help.GetCurrentPath()
	if err != nil {
		log.Fatalln("path can`t read")
		return ""
	}
	AllLogDir := currPath + logDir
	if _, err := os.Stat(AllLogDir); err != nil {
		err = os.MkdirAll(AllLogDir, 0711)
		if err != nil {
			log.Fatalln("create Dir ERROR: ", err)
		}
	}
	return AllLogDir
}

func structLog(level string) *log.Logger {
	var w io.Writer
	if debuger {
		w = io.MultiWriter(os.Stdout, getFile(level))
	} else {
		w = io.MultiWriter(getFile(level))
	}
	return log.New(w, level+": ", logInfo)
}

func getLog(level string) *log.Logger {
	switch level {
	case "Info":
		if infoLog == nil {
			infoLog = structLog(level)
		}
		return infoLog
	case "Warning":
		if warnLog == nil {
			warnLog = structLog(level)
		}
		return warnLog
	case "Error":
		if errLog == nil {
			errLog = structLog(level)
		}
		return errLog
	case "Fatal":
		if fatalLog == nil {
			fatalLog = structLog(level)
		}
		return fatalLog
	default:
		if infoLog == nil {
			infoLog = log.New(io.MultiWriter(os.Stdout, getFile("Info")), "Info: ", logInfo)
		}
		return infoLog
	}
}

/**
 * 写入
 * @param   loger  *log.Logger      [description]
 * @param   args ...interface{} [description]
 */
func write(loger *log.Logger, args ...interface{}) {
	logStrs := make([]string, len(args))
	for i, _ := range args {
		logStrs[i] = fmt.Sprintf("%+v", args[i])
	}
	if v := reflect.ValueOf(loger).MethodByName("Output"); v.String() != "<invalid Value>" {
		v.Call([]reflect.Value{
			reflect.ValueOf(6),
			reflect.ValueOf(strings.Join(logStrs, " ") + "\n"),
		})
	}
}
