package logx

import (
	"fmt"
	"runtime"
	"strings"
)

// CallerPrettyfier is func for format file name
func CallerPrettyfier(f *runtime.Frame) (string, string) {
	i := strings.LastIndex(f.File, "/")
	if i >= 0 {
		f.File = f.File[i+1:]
	}

	i = strings.LastIndex(f.Function, "/")
	if i >= 0 {
		f.Function = f.Function[i+1:]
	}

	return f.Function, fmt.Sprintf("%s:%d", f.File, f.Line)
}
