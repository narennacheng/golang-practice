package geek

import (
	"os"
	"syscall"
)

var (
	ShutdowmSignals = []os.Signal{
		os.Interrupt, os.Kill, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL,
		syscall.SIGABRT, syscall.SIGTERM, syscall.SIGSYS, syscall.SIGSTOP,
	}
)
