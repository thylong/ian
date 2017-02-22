package command

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// ExecuteCommand a command and print output from stdout.
func ExecuteCommand(subCmd *exec.Cmd) (err error) {
	cmdReader, err := subCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error creating StdoutPipe for Cmd: %v StdErr: %v", err, os.Stderr)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = subCmd.Start()
	if err != nil {
		return fmt.Errorf("Error starting Cmd: %v StdErr: %v", err, os.Stderr)
	}

	err = subCmd.Wait()
	if err != nil {
		return fmt.Errorf("Error waiting for Cmd: %v StdErr: %v", err, os.Stderr)
	}

	return nil
}
