// +build darwin freebsd linux netbsd openbsd

package clicontext

import "syscall"

func init() {
	defaultInterruptSignals = append(defaultInterruptSignals, syscall.SIGTERM)
}
