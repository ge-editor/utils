package utils

import (
	"bytes"
	"os/exec"
)

// ExecShellCommand はシェル経由でコマンドを実行し stdout を返す
func ExecShellCommand(cmd string) (string, error) {
	c := exec.Command("sh", "-c", cmd)

	var out bytes.Buffer
	var errOut bytes.Buffer
	c.Stdout = &out
	c.Stderr = &errOut

	err := c.Run()
	if err != nil {
		return errOut.String(), err
	}
	return out.String(), nil
}
