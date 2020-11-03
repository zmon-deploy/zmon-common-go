package system

import (
	"bytes"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
	"time"
)

func ExecuteCmd(c *exec.Cmd, timeout time.Duration) (string, error) {
	var outBuf bytes.Buffer
	c.Stdout = &outBuf
	c.Stderr = &outBuf

	if err := c.Start(); err != nil {
		return "", errors.Wrap(err, "failed to start cmd")
	}

	err := waitTimeout(c, timeout)
	if err != nil {
		return "", errors.Wrap(err, "failed to wait for process termination")
	}

	return strings.TrimSpace(string(outBuf.Bytes())), nil
}

func waitTimeout(c *exec.Cmd, timeout time.Duration) error {
	var err error
	isDone := false

	timer := time.AfterFunc(timeout, func() {
		if isDone == true {
			return
		}
		err = c.Process.Kill()
	})

	_ = c.Wait()
	isDone = true
	timer.Stop()

	return err
}
