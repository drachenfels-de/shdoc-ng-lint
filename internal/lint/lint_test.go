package lint_test

import (
	"testing"

	"github.com/drachenfels-de/shdoc-ng-lint/internal/lint"
	shdoc "github.com/jdevera/shdoc-ng/shdoc"
)

func TestDocument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		src   string
		opts  lint.Options
		want  []lint.Issue
		parse shdoc.ParseOptions
	}{
		{
			name: "passes when all required tags are present",
			src: `# @description Say hello
# @example
#   hello_world "Jane"
# @arg $1 String name to greet
# @exitcode 0 Printed greeting
hello_world() {
  echo "hello $1"
}
`,
		},
		{
			name: "fails missing description example args and exitcode",
			src: `hello_world() {
  echo "hello"
}
`,
			parse: shdoc.ParseOptions{IncludeUndocumented: true},
			want: []lint.Issue{
				{FunctionName: "hello_world", Message: "missing non-empty @description"},
				{FunctionName: "hello_world", Message: "missing @example"},
				{FunctionName: "hello_world", Message: "must define @noargs or at least one non-zero @arg"},
				{FunctionName: "hello_world", Message: "missing at least one @exitcode"},
			},
		},
		{
			name: "allows noargs",
			src: `# @description Say hello
# @example
#   hello_world
# @noargs
# @exitcode 0 Printed greeting
hello_world() {
  echo "hello"
}
`,
		},
		{
			name: "rejects zero arg only",
			src: `# @description Bootstrap shell
# @example
#   bootstrap "$0"
# @arg $0 String current script path
# @exitcode 0 Printed script path
bootstrap() {
  printf '%s\n' "$0"
}
`,
			want: []lint.Issue{
				{FunctionName: "bootstrap", Message: "must define @noargs or at least one non-zero @arg"},
			},
		},
		{
			name: "enforces prefix when configured",
			src: `# @description Say hello
# @example
#   hello_world "Jane"
# @arg $1 String name to greet
# @exitcode 0 Printed greeting
hello_world() {
  echo "hello $1"
}
`,
			opts: lint.Options{FunctionPrefix: "doc_"},
			want: []lint.Issue{
				{FunctionName: "hello_world", Message: "must start with prefix \"doc_\""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doc, warns := shdoc.ParseDocumentWithOptions(tt.src, tt.parse)
			if len(warns) != 0 {
				t.Fatalf("ParseDocumentWithOptions() warnings = %v", warns)
			}

			got := lint.Document(doc, tt.opts)
			if len(got) != len(tt.want) {
				t.Fatalf("Document() issue count = %d, want %d; issues=%v", len(got), len(tt.want), got)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("Document() issue[%d] = %#v, want %#v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
