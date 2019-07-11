package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os/exec"
)

// read to buffer
func readToBuffer(buffer bytes.Buffer, rd io.Reader) {
	reader := bufio.NewReader(rd)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		buffer.WriteString(line)
	}
}

// ExecCommand execute command, return result and error
func ExecCommand(commandName string, params []string) (*exec.Cmd, string, error) {
	cmd := exec.Command(commandName, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, "", err
	}

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Start()
	readToBuffer(outBuffer, stdout)
	readToBuffer(errBuffer, stderr)
	cmd.Wait()
	errstr := errBuffer.String()
	if len(errstr) > 0 {
		return cmd, outBuffer.String(), errors.New(errBuffer.String())
	}
	return cmd, outBuffer.String(), nil
}
