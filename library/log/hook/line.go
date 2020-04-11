package hook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// line number hook for log the call context,
type LineHook struct {
	// skip为遍历调用栈开始的索引位置
	Skip   int
	levels []logrus.Level
}

// Levels implement levels
func (h *LineHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (h *LineHook) Fire(entry *logrus.Entry) error {
	entry.Data["stack"] = h.findCaller(h.Skip)
	return nil
}

func (h *LineHook) findCaller(skip int) string {
	file := ""
	line := 0
	var pc uintptr
	// 遍历调用栈的最大索引为第15层.
	for i := 0; i < 15; i++ {
		file, line, pc = h.getCaller(skip + i)
		// 过滤无用的栈
		if !strings.HasPrefix(file, "hook/") &&
			!strings.HasPrefix(file, "logrus") &&
			!strings.HasPrefix(file, "log/") {
			break
		}
	}
	fullFnName := runtime.FuncForPC(pc)
	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		// 取得函数名
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}
	return fmt.Sprintf("%s:%d:%s()", file, line, fnName)
}

func (h *LineHook) getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0
	// 获取包名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}
