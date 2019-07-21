package exec

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	tests := gentests(false)
	for _, tt := range tests {
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		cmd := Command(tt.cmd[0], tt.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if strings.TrimSuffix(stdout.String(), "\n") != tt.want {
			t.Errorf("%s: want = %#v, got = %#v", tt.name, tt.want, stdout.String())
		}
		_, err = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep", tt.want)).Output()
		if err == nil {
			t.Errorf("%s", "the process has not exited")
		}
	}
}

func TestCommandContext(t *testing.T) {
	tests := gentests(false)
	for _, tt := range tests {
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		ctx := context.Background()
		cmd := CommandContext(ctx, tt.cmd[0], tt.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if strings.TrimSuffix(stdout.String(), "\n") != tt.want {
			t.Errorf("%s: want = %#v, got = %#v", tt.name, tt.want, stdout.String())
		}
		_, err = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep", tt.want)).Output()
		if err == nil {
			t.Errorf("%s", "the process has not exited")
		}
	}
}

func TestCommandContextCancel(t *testing.T) {
	tests := gentests(true)
	for _, tt := range tests {
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		cmd := CommandContext(ctx, tt.cmd[0], tt.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Start()
		if err != nil {
			t.Fatalf("%v", err)
		}
		time.Sleep(100 * time.Millisecond)
		cancel()
		_, err = exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep", tt.want)).Output()
		if err == nil {
			t.Errorf("%s", "the process has not exited")
		}
	}
}

type testcase struct {
	name string
	cmd  []string
	want string
}

func gentests(withSleep bool) []testcase {
	tests := []testcase{}
	r := random()
	tests = append(tests, testcase{"simple echo", []string{"echo", r}, r})
	r = random()
	tests = append(tests, testcase{"bash -c echo", []string{"bash", "-c", fmt.Sprintf("echo %s", r)}, r})
	if withSleep {
		r = "123456"
		tests = append(tests, testcase{"sleep", []string{"sleep", r}, r})
		r = "654321"
		tests = append(tests, testcase{"bash -c sleep", []string{"bash", "-c", fmt.Sprintf("sleep %s", r)}, r})
	}
	return tests
}

func random() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}
