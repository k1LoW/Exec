package exec

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// Exec represents an command executer
type Exec struct {
	Signal          os.Signal
	KillAfterCancel time.Duration // TODO
}

// CommandContext returns *os/exec.Cmd with Setpgid = true
// When ctx cancelled, `github.com/k1LoW/exec.CommandContext` send signal to process group
func (e *Exec) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := command(name, arg...)
	go func() {
		select {
		case <-ctx.Done():
			err := terminate(cmd, e.Signal)
			if err != nil {
				// :thinking:
				return
			}
		}
	}()
	return cmd
}

// LookPath is os/exec.LookPath
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Command returns *os/exec.Cmd with Setpgid = true
func Command(name string, arg ...string) *exec.Cmd {
	return command(name, arg...)
}

// CommandContext returns *os/exec.Cmd with Setpgid = true
// When ctx cancelled, `github.com/k1LoW/exec.CommandContext` send signal to process group
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	e := &Exec{
		Signal:          os.Kill,
		KillAfterCancel: -1,
	}
	return e.CommandContext(ctx, name, arg...)
}

// TerminateCommand send signal to cmd.Process.Pid process group
func TerminateCommand(cmd *exec.Cmd, sig os.Signal) error {
	return terminate(cmd, sig)
}

// KillCommand send syscall.SIGKILL to cmd.Process.Pid process group
func KillCommand(cmd *exec.Cmd) error {
	return killall(cmd)
}
