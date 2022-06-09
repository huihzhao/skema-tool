package console

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func Info(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func Infof(format string, args ...any) {
	fmt.Printf(format, args...)
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf("[ERROR]"+format, args...)
}

func Fatalf(format string, args ...any) {
	log.Fatalf(format, args...)
}

func ExecCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	return executeCmd(cmd)
}

func ExecCommandWithPath(path string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = path
	return executeCmd(cmd)
}

func FatalIfError(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			Fatalf(err.Error())
		}

		Fatalf("%s\n%s\n", msg, err.Error())
	}
}

func executeCmd(cmd *exec.Cmd) error {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			Errorf(scanner.Text())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			Infof(scanner.Text())
		}
	}()
	err := cmd.Wait()
	return err
}
