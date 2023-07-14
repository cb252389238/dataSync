//go:build windows
// +build windows

package tools

import (
	"os"
	"os/exec"
)

func Cls() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
