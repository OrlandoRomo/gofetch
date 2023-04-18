package macos

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var isBootTimeCommand = true

func TestUptimeHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_UPTIME_BOOTTIME") != "1" && os.Getenv("GO_WANT_HELPER_PROCESS_UPTIME") != "1" && os.Getenv("GO_WANT_HELPER_PROCESS_FAILURE") != "1" {
		return
	}
	if os.Getenv("GO_WANT_HELPER_PROCESS_UPTIME_BOOTTIME") == "1" {
		fmt.Fprintf(os.Stdout, "{ sec = 1656944602, usec = 201499 } Mon Jul  4 09:23:22 2022")
	}

	if os.Getenv("GO_WANT_HELPER_PROCESS_UPTIME") == "1" {
		fmt.Fprintf(os.Stdout, "123456")
	}

	if os.Getenv("GO_WANT_HELPER_PROCESS_FAILURE") == "1" {
		os.Exit(1)
	}

	os.Exit(0)
}
func TestGetUptime(t *testing.T) {
	tcs := []struct {
		Desc            string
		Expected        string
		EnvCommands     []string
		FakeExecCommand func(envs []string) func(command string, arg ...string) *exec.Cmd
	}{
		{
			Desc:     "success - received uptime",
			Expected: "1 day(s), 10 hour(s), 17 minute(s)",
			EnvCommands: []string{
				"GO_WANT_HELPER_PROCESS_UPTIME_BOOTTIME=1",
				"GO_WANT_HELPER_PROCESS_UPTIME=1",
			},
			FakeExecCommand: func(envs []string) func(command string, arg ...string) *exec.Cmd {
				return func(command string, args ...string) *exec.Cmd {
					cs := []string{"-test.run=TestUptimeHelper", "--", command}
					cs = append(cs, args...)
					cmd := exec.Command(os.Args[0], cs...)
					cmd.Env = envs
					envs = envs[1:]
					return cmd
				}
			},
		},
		{
			Desc:     "failure - unable to get boot time",
			Expected: "Unknown",
			EnvCommands: []string{
				"GO_WANT_HELPER_PROCESS_FAILURE=1",
			},
			FakeExecCommand: func(envs []string) func(command string, arg ...string) *exec.Cmd {
				return func(command string, args ...string) *exec.Cmd {
					cs := []string{"-test.run=TestUptimeHelper", "--", command}
					cs = append(cs, args...)
					cmd := exec.Command(os.Args[0], cs...)
					cmd.Env = envs
					envs = envs[1:]
					return cmd
				}
			},
		},
		{
			Desc:     "failure - unable to get uptime in seconds",
			Expected: "Unknown",
			EnvCommands: []string{
				"GO_WANT_HELPER_PROCESS_UPTIME_BOOTTIME=1",
				"GO_WANT_HELPER_PROCESS_FAILURE=1",
			},
			FakeExecCommand: func(envs []string) func(command string, arg ...string) *exec.Cmd {
				return func(command string, args ...string) *exec.Cmd {
					cs := []string{"-test.run=TestUptimeHelper", "--", command}
					cs = append(cs, args...)
					cmd := exec.Command(os.Args[0], cs...)
					cmd.Env = envs
					envs = envs[1:]
					return cmd
				}
			},
		},
	}

	for _, tt := range tcs {
		t.Run(tt.Desc, func(t *testing.T) {
			execCommand = tt.FakeExecCommand(tt.EnvCommands)
			mac := New()
			uptime := mac.GetUptime()
			if uptime != tt.Expected {
				t.Fatalf("received %s but expected %s", uptime, tt.Expected)
			}
		})
	}
}
