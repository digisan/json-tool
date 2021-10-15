package jsontool

import (
	"reflect"
	"testing"
)

func TestIndexAll(t *testing.T) {
	type args struct {
		s   string
		sub string
	}
	tests := []struct {
		name       string
		args       args
		wantStarts []int
		wantEnds   []int
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				s:   "abcdeabcdabe  ab",
				sub: "ab",
			},
			wantStarts: []int{0, 5, 9, 14},
			wantEnds:   []int{2, 7, 11, 16},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStarts, gotEnds := IndexAll(tt.args.s, tt.args.sub)
			if !reflect.DeepEqual(gotStarts, tt.wantStarts) {
				t.Errorf("IndexAll() gotStarts = %v, want %v", gotStarts, tt.wantStarts)
			}
			if !reflect.DeepEqual(gotEnds, tt.wantEnds) {
				t.Errorf("IndexAll() gotEnds = %v, want %v", gotEnds, tt.wantEnds)
			}
		})
	}
}
