package misc

import (
	"fmt"
	"runtime"
	"strings"
)

type errorable interface {
	Error(...interface{})
}

func LogPanic(logger errorable) {
	if r := recover(); r != nil {
		messages := []string{
			fmt.Sprintf("panic: %s", r),
		}

		i := 0
		fnName, file, line, ok := runtime.Caller(i)
		for ok {
			message := fmt.Sprintf("[fn]: %s, [file]: %s, [line]: %d", runtime.FuncForPC(fnName).Name(), file, line)
			messages = append(messages, message)
			i++
			fnName, file, line, ok = runtime.Caller(i)
		}

		logger.Error(strings.Join(messages, "\n"))
	}
}
