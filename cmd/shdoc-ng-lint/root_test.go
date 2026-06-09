package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunLintSuccess(t *testing.T) {
	t.Parallel()

	stdin := strings.NewReader(`# @description Say hello
# @example
#   hello_world "Jane"
# @arg $1 String name to greet
# @exitcode 0 Printed greeting
hello_world() {
  echo "hello $1"
}
`)

	origStdin := osStdin
	osStdin = stdin
	t.Cleanup(func() {
		osStdin = origStdin
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := runLint(&stdout, &stderr, &config{})
	if err != nil {
		t.Fatalf("runLint() error = %v", err)
	}

	if got := stdout.String(); got != "lint passed\n" {
		t.Fatalf("stdout = %q, want %q", got, "lint passed\n")
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunLintFailure(t *testing.T) {
	t.Parallel()

	stdin := strings.NewReader("hello_world() {\n  echo hello\n}\n")

	origStdin := osStdin
	osStdin = stdin
	t.Cleanup(func() {
		osStdin = origStdin
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := runLint(&stdout, &stderr, &config{includeUndocumented: true})
	if err == nil {
		t.Fatal("runLint() error = nil, want failure")
	}

	exitErr, ok := err.(interface{ ExitCode() int })
	if !ok {
		t.Fatalf("runLint() error type = %T, want ExitCode", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("ExitCode() = %d, want 1", exitErr.ExitCode())
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	if !strings.Contains(stderr.String(), `function "hello_world": missing non-empty @description`) {
		t.Fatalf("stderr = %q, want lint issue", stderr.String())
	}
}

func TestRunLintWerror(t *testing.T) {
	t.Parallel()

	stdin := strings.NewReader(`# @description Say hello
# @warning
# @example
#   hello_world
# @noargs
# @exitcode 0 Printed greeting
hello_world() {
  echo "hello"
}
`)

	origStdin := osStdin
	osStdin = stdin
	t.Cleanup(func() {
		osStdin = origStdin
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := runLint(&stdout, &stderr, &config{werror: true})
	if err == nil {
		t.Fatal("runLint() error = nil, want failure")
	}

	exitErr, ok := err.(interface{ ExitCode() int })
	if !ok {
		t.Fatalf("runLint() error type = %T, want ExitCode", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("ExitCode() = %d, want 1", exitErr.ExitCode())
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	if !strings.Contains(stderr.String(), "warning: Empty value: @warning requires a message") {
		t.Fatalf("stderr = %q, want parser warning", stderr.String())
	}
}
