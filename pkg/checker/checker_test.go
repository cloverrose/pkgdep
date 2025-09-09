package checker

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/cloverrose/pkgdep/pkg/orderedmap"
)

func TestChecker_IsAllowedDependency(t *testing.T) {
	var (
		simpleDep = func() orderedmap.OrderedMap {
			om := orderedmap.New()
			om.Set("a", []string{"b"})
			om.Set("b", []string{"c"})
			return *om
		}

		regexpDep = func() orderedmap.OrderedMap {
			om := orderedmap.New()
			om.Set("^a$", []string{"^b$"})
			om.Set("^a+$", []string{"^b+$"}) // This relaxed one matches only when above rule does not match.
			return *om
		}

		patternDep = func() orderedmap.OrderedMap {
			om := orderedmap.New()
			om.Set(`^(?P<ax>a+)(?P<bx>b+)(?P<cx>c+)$`, []string{`^{{ .ax }}/{{ .bx }}/{{ .cx }}$`})
			return *om
		}
	)

	tests := []struct {
		name         string
		dependencies func() orderedmap.OrderedMap
		recorder     func(ctrl *gomock.Controller) recorder
		from         string
		to           string
		want         bool
	}{
		{
			name:         "simple match a to b",
			dependencies: simpleDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage("a", "b")
				return m
			},
			from: "a",
			to:   "b",
			want: true,
		},
		{
			name:         "simple match b to c",
			dependencies: simpleDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage("b", "c")
				return m
			},
			from: "b",
			to:   "c",
			want: true,
		},
		{
			name:         "simple unmatch",
			dependencies: simpleDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "c",
			to:   "d",
			want: false,
		},
		{
			name:         "regexp match",
			dependencies: regexpDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage("^a$", "^b$")
				return m
			},
			from: "a",
			to:   "b",
			want: true,
		},
		{
			name:         "regexp match relax",
			dependencies: regexpDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage("^a+$", "^b+$")
				return m
			},
			from: "aaa",
			to:   "bb",
			want: true,
		},
		{
			name:         "regexp unmatch",
			dependencies: regexpDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "b",
			to:   "a",
			want: false,
		},
		{
			name:         "pattern match",
			dependencies: patternDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage(`^(?P<ax>a+)(?P<bx>b+)(?P<cx>c+)$`, `^{{ .ax }}/{{ .bx }}/{{ .cx }}$`)
				return m
			},
			from: "abbccc",
			to:   "a/bb/ccc",
			want: true,
		},
		{
			name:         "pattern unmatch",
			dependencies: patternDep,
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "abc",
			to:   "aa/bb/cc",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			checker := New(tt.dependencies(), tt.recorder(ctrl))
			got := checker.IsAllowedDependency(tt.from, tt.to)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("IsAllowedDependency() = (-want +got):\n%s", diff)
			}
		})
	}
}
