package tidy

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTidy_tidyCore(t *testing.T) {
	tests := []struct {
		name         string
		configInput  string
		inspectInput string
		want         string
	}{
		{
			name: "remove single unused rule",
			configInput: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
    - a/b/d
  a/c:
    - a/c/d
`,
			inspectInput: `a/b,a/b/c
a/c,a/c/d`,
			want: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
  a/c:
    - a/c/d
`,
		},
		{
			name: "remove empty key",
			configInput: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
  a/c:
    - a/c/d
`,
			inspectInput: `a/c,a/c/d`,
			want: `targetPackagePrefixList:
  - a
dependencies:
  a/c:
    - a/c/d
`,
		},
		{
			name: "no rules to remove",
			configInput: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
  a/c:
    - a/c/d
`,
			inspectInput: `a/b,a/b/c
a/c,a/c/d`,
			want: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
  a/c:
    - a/c/d
`,
		},
		{
			name: "remove all rules",
			configInput: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
  a/c:
    - a/c/d
`,
			inspectInput: ``,
			want: `targetPackagePrefixList:
  - a
dependencies: {}
`,
		},
		{
			name: "remove multiple rules from same key",
			configInput: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/c
    - a/b/d
    - a/b/e
  a/c:
    - a/c/d
`,
			inspectInput: `a/b,a/b/d
a/c,a/c/d`,
			want: `targetPackagePrefixList:
  - a
dependencies:
  a/b:
    - a/b/d
  a/c:
    - a/c/d
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			configReader := strings.NewReader(tt.configInput)
			configReaderForDeps := strings.NewReader(tt.configInput)
			inspectReader := strings.NewReader(tt.inspectInput)
			out := &bytes.Buffer{}

			p := &Tidy{}
			if err := p.tidyCore(configReader, configReaderForDeps, inspectReader, out); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := out.String()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("tidyCore() (-want,+got):\n%s", diff)
			}
		})
	}
}
