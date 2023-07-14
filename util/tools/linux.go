//go:build linux
// +build linux

package tools

import (
	"os"
	"os/exec"
)

func Cls() {
	cmd := exec.Command("/bin/sh", "-c", "clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
