package geek

import (
	"os"
	"syscall"
)

var (
	ShutdowmSignalsWin = []os.Signal{
		os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGILL,
		syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGTERM, syscall.SIGABRT,
		os.Kill, syscall.SIGKILL, syscall.SIGQUIT,
	}
)
