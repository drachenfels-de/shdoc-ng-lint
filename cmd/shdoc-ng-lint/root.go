package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/drachenfels-de/shdoc-ng-lint/internal/lint"
	shdoc "github.com/jdevera/shdoc-ng/shdoc"
	"github.com/spf13/cobra"
)

var osStdin io.Reader = os.Stdin

type config struct {
	inputs              []string
	functionPrefix      string
	includeUndocumented bool
	includeInternal     bool
	includeAll          bool
	werror              bool
}

func newRootCommand(stdout io.Writer, stderr io.Writer) *cobra.Command {
	cfg := &config{}

	cmd := &cobra.Command{
		Use:     "shdoc-ng-lint [flags] [file ...]",
		Short:   "Lint shdoc-ng annotated shell script documentation",
		Version: version,
		Long: `Lint shell script documentation parsed by shdoc-ng.

By default the command reads from stdin when no file arguments are provided.`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.inputs = args
			return runLint(stdout, stderr, cfg)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().StringVar(&cfg.functionPrefix, "function-prefix", "", "Require every function name to start with this prefix")
	cmd.Flags().BoolVar(&cfg.includeUndocumented, "include-undocumented", false, "Include undocumented functions in the parsed document")
	cmd.Flags().BoolVar(&cfg.includeInternal, "include-internal", false, "Include @internal functions in the parsed document")
	cmd.Flags().BoolVar(&cfg.includeAll, "include-all", false, "Shorthand for --include-undocumented --include-internal")
	cmd.Flags().BoolVar(&cfg.werror, "werror", false, "Treat all warnings as linting errors")

	return cmd
}

func run() int {
	cmd := newRootCommand(os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		var exitErr interface{ ExitCode() int }
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func runLint(stdout io.Writer, stderr io.Writer, cfg *config) error {
	targets := cfg.inputs
	if len(targets) == 0 {
		targets = []string{"-"}
	}

	totalIssues := 0
	totalWarns := 0
	for _, target := range targets {
		src, displayName, err := readInput(target)
		if err != nil {
			return err
		}

		doc, warns := shdoc.ParseDocumentWithOptions(src, shdoc.ParseOptions{
			IncludeUndocumented: cfg.includeUndocumented || cfg.includeAll,
			IncludeInternal:     cfg.includeInternal || cfg.includeAll,
		})

		for _, warning := range warns {
			if _, err := fmt.Fprintf(stderr, "%s:%d:%d: warning: %s\n", displayName, warning.Line, warning.Col+1, warning.Message); err != nil {
				return fmt.Errorf("writing warning: %w", err)
			}
		}

		totalWarns += len(warns)

		issues := lint.Document(doc, lint.Options{FunctionPrefix: cfg.functionPrefix})
		for _, issue := range issues {
			if _, err := fmt.Fprintf(stderr, "%s: function %q: %s\n", displayName, issue.FunctionName, issue.Message); err != nil {
				return fmt.Errorf("writing issue: %w", err)
			}
		}

		totalIssues += len(issues)
	}

	if totalIssues > 0 || (cfg.werror && totalWarns > 0) {
		return lint.ExitError{Code: 1, Message: fmt.Sprintf("lint failed with %d issue(s) and %d warning(s)", totalIssues, totalWarns)}
	}

	if _, err := fmt.Fprintln(stdout, "lint passed"); err != nil {
		return fmt.Errorf("writing result: %w", err)
	}

	return nil
}

func readInput(target string) (string, string, error) {
	if target == "-" {
		src, err := io.ReadAll(osStdin)
		if err != nil {
			return "", "", fmt.Errorf("reading stdin: %w", err)
		}

		return string(src), "<stdin>", nil
	}

	src, err := os.ReadFile(target)
	if err != nil {
		return "", "", fmt.Errorf("reading %s: %w", target, err)
	}

	return string(src), target, nil
}
