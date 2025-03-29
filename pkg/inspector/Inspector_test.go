package inspector

import (
	"testing"
)

func TestInspector_IsRecorded(t *testing.T) {
	t.Parallel()

	type fields struct {
		usedRules map[string][]string
	}
	type args struct {
		frm string
		to  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "nil",
			fields: fields{
				usedRules: nil,
			},
			args: args{
				frm: "from1",
				to:  "to1",
			},
			want: false,
		},
		{
			name: "empty usedRules",
			fields: fields{
				usedRules: map[string][]string{},
			},
			args: args{
				frm: "from1",
				to:  "to1",
			},
			want: false,
		},
		{
			name: "not match",
			fields: fields{
				usedRules: map[string][]string{
					"from2": {"to2"},
				},
			},
			args: args{
				frm: "from1",
				to:  "to1",
			},
			want: false,
		},
		{
			name: "from match but to not match",
			fields: fields{
				usedRules: map[string][]string{
					"from1": {"to2"},
				},
			},
			args: args{
				frm: "from1",
				to:  "to1",
			},
			want: false,
		},
		{
			name: "match",
			fields: fields{
				usedRules: map[string][]string{
					"from1": {"to1"},
				},
			},
			args: args{
				frm: "from1",
				to:  "to1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			i := &Inspector{
				usedRules: tt.fields.usedRules,
			}
			if got := i.IsRecorded(tt.args.frm, tt.args.to); got != tt.want {
				t.Errorf("IsRecorded() = %v, want %v", got, tt.want)
			}
		})
	}
}
