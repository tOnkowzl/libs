package logx

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_CallerPrettyfier(b *testing.B) {
	f := &runtime.Frame{
		File:     "ci.inno.ktb/microservices/go/logs.git/v3/main.go",
		Function: "ci.inno.ktb/microservices/go/logs.git/v3/middlewarex.LoggerWith.func1.1",
		Line:     10,
	}

	for i := 1; i <= b.N; i++ {
		CallerPrettyfier(f)
	}
}

func Test_CallerPrettyfier(t *testing.T) {
	f := &runtime.Frame{
		File:     "ci.inno.ktb/microservices/go/logs.git/v3/main.go",
		Function: "ci.inno.ktb/microservices/go/logs.git/v3/middlewarex.LoggerWith.func1.1",
		Line:     10,
	}

	fl, fn := CallerPrettyfier(f)

	assert.Equal(t, fl, "middlewarex.LoggerWith.func1.1")
	assert.Equal(t, fn, "main.go:10")
}
