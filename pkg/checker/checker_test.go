package checker

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"

	"github.com/cloverrose/pkgdep/pkg/orderedmap"
)

func TestChecker_CheckDependency(t *testing.T) {
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

		globalDataDep = func() orderedmap.OrderedMap {
			om := orderedmap.New()
			om.Set(`^(?P<moduleName>(a|b|c))$`, []string{
				`^{{ index .globalData .moduleName }}$`,
				`^{{ index .globalData "module" .moduleName }}$`,
			})
			return *om
		}

		complexGlobalDataDep = func() orderedmap.OrderedMap {
			om := orderedmap.New()
			om.Set(`^(?P<moduleName>[^/]+)/(?P<layerName>[^/]+)$`, []string{
				`^{{ .moduleName }}/{{ index .globalData .moduleName .layerName }}$`,
			})
			return *om
		}

		complexGlobalData = func() map[string]any {
			return map[string]any{
				// imagine module a and b have application, domain, infra layer.
				"a": map[string]any{
					"infra":       "(infra)",
					"domain":      "(domain|infra)",
					"application": "(application|domain|infra)",
				},
				"b": map[string]any{
					"infra":       "(infra)",
					"domain":      "(domain|infra)",
					"application": "(application|domain|infra)",
				},
				// and module c has only application and infra layer.
				"c": map[string]any{
					"infra":       "(infra)",
					"application": "(application|infra)",
				},
			}
		}
	)

	tests := []struct {
		name         string
		dependencies func() orderedmap.OrderedMap
		globalData   map[string]any
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
		{
			name: "pattern match handle missing key",
			dependencies: func() orderedmap.OrderedMap {
				om := orderedmap.New()
				om.Set(`^(?P<ax>a+)$`, []string{`^{{ .badkey }}$`})
				return *om
			},
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "aaa",
			to:   "foo",
			want: false,
		},
		{
			name:         "globalData match",
			dependencies: globalDataDep,
			globalData: map[string]any{
				"a": "(a|b|c)",
				"b": "(b|c)",
				"c": "(c)",
			},
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage(`^(?P<moduleName>(a|b|c))$`, `^{{ index .globalData .moduleName }}$`)
				return m
			},
			from: "a",
			to:   "b",
			want: true,
		},
		{
			name:         "globalData nested static key match",
			dependencies: globalDataDep,
			globalData: map[string]any{
				"module": map[string]any{
					"a": "(a|b|c)",
					"b": "(b|c)",
					"c": "(c)",
				},
			},
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage(`^(?P<moduleName>(a|b|c))$`, `^{{ index .globalData "module" .moduleName }}$`)
				return m
			},
			from: "b",
			to:   "c",
			want: true,
		},
		{
			name:         "complex globalData inside module a (match)",
			dependencies: complexGlobalDataDep,
			globalData:   complexGlobalData(),
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage(gomock.Any(), gomock.Any())
				return m
			},
			from: "a/domain",
			to:   "a/infra",
			want: true,
		},
		{
			name:         "complex globalData cross module (unmatch)",
			dependencies: complexGlobalDataDep,
			globalData:   complexGlobalData(),
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "a/application",
			to:   "b/domain",
			want: false,
		},
		{
			name:         "complex globalData inside module c (match)",
			dependencies: complexGlobalDataDep,
			globalData:   complexGlobalData(),
			recorder: func(ctrl *gomock.Controller) recorder {
				m := NewMockrecorder(ctrl)
				m.EXPECT().RecordUsage(gomock.Any(), gomock.Any())
				return m
			},
			from: "c/application",
			to:   "c/infra",
			want: true,
		},
		{
			name:         "complex globalData referring invalid layer (unmatch)",
			dependencies: complexGlobalDataDep,
			globalData:   complexGlobalData(),
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "c/application",
			to:   "c/domain", // module c does not have domain layer
			want: false,
		},
		{
			name: "globalData key conflicts",
			dependencies: func() orderedmap.OrderedMap {
				om := orderedmap.New()
				// globalData is used for pattern name.
				om.Set(`^(?P<globalData>a+)$`, []string{`^{{ index .globalData "hello" }}$`})
				return *om
			},
			globalData: map[string]any{
				"hello": "world",
			},
			recorder: func(ctrl *gomock.Controller) recorder {
				return NewMockrecorder(ctrl)
			},
			from: "aaa",
			to:   "world",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Discard noisy log
			slog.SetDefault(slog.New(slog.DiscardHandler))

			ctrl := gomock.NewController(t)
			checker := New(tt.dependencies(), tt.globalData, tt.recorder(ctrl))
			got := checker.CheckDependency(tt.from, tt.to)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CheckDependency() = (-want +got):\n%s", diff)
			}
		})
	}
}
