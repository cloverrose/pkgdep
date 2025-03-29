package orderedmap

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrderedMap_Iter(t *testing.T) {
	t.Parallel()
	type fields struct {
		Keys   []string
		Values map[string][]string
	}
	type pair struct {
		Key   string
		Value []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []pair
	}{
		{
			name: "nil",
			fields: fields{
				Keys:   nil,
				Values: nil,
			},
			want: nil,
		},
		{
			name: "empty",
			fields: fields{
				Keys:   []string{},
				Values: map[string][]string{},
			},
			want: nil,
		},
		{
			name: "one",
			fields: fields{
				Keys:   []string{"a"},
				Values: map[string][]string{"a": {"1", "2", "3"}},
			},
			want: []pair{{Key: "a", Value: []string{"1", "2", "3"}}},
		},
		{
			name: "two",
			fields: fields{
				Keys:   []string{"b", "a"},
				Values: map[string][]string{"a": {"1", "2", "3"}, "b": {"4", "5", "6"}},
			},
			want: []pair{
				{Key: "b", Value: []string{"4", "5", "6"}},
				{Key: "a", Value: []string{"1", "2", "3"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			om := &OrderedMap{
				keys:   tt.fields.Keys,
				values: tt.fields.Values,
			}
			var got []pair
			for k, v := range om.Iter() {
				got = append(got, pair{Key: k, Value: v})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("OrderedMap.Iter() (-want,+got):\n%s", diff)
			}
		})
	}
}
