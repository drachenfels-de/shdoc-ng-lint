package lint

import (
	"fmt"
	"strings"

	shdoc "github.com/jdevera/shdoc-ng/shdoc"
)

type Options struct {
	FunctionPrefix string
}

type Issue struct {
	FunctionName string
	Message      string
}

type ExitError struct {
	Code    int
	Message string
}

func (e ExitError) Error() string {
	return e.Message
}

func (e ExitError) ExitCode() int {
	return e.Code
}

func Document(doc shdoc.Document, opts Options) []Issue {
	var issues []Issue

	for _, fn := range doc.AllFunctions() {
		if strings.TrimSpace(fn.Description) == "" {
			issues = append(issues, issue(fn.Name, "missing non-empty @description"))
		}

		if strings.TrimSpace(fn.Example) == "" {
			issues = append(issues, issue(fn.Name, "missing non-empty @example"))
		}

		if !fn.IsNoArgs && !hasNonZeroArg(fn.Args) {
			issues = append(issues, issue(fn.Name, "must define @noargs or at least one non-zero @arg"))
		}

		if len(fn.ExitCodes) == 0 {
			issues = append(issues, issue(fn.Name, "missing at least one @exitcode"))
		}

		if opts.FunctionPrefix != "" && !strings.HasPrefix(fn.Name, opts.FunctionPrefix) {
			issues = append(issues, issue(fn.Name, fmt.Sprintf("must start with prefix %q", opts.FunctionPrefix)))
		}
	}

	return issues
}

func hasNonZeroArg(args []shdoc.Arg) bool {
	for _, arg := range args {
		if strings.TrimSpace(arg.Name) == "$0" {
			continue
		}

		if strings.TrimSpace(arg.Name) != "" {
			return true
		}
	}

	return false
}

func issue(functionName string, message string) Issue {
	return Issue{
		FunctionName: functionName,
		Message:      message,
	}
}
