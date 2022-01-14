package compiler

import (
	"testing"
)

func TestNewReferencePath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"$", false},
		{"$.aaa", false},
		{"$.ledgers.branch[0].pending.count", false},
		{"$.ledgers.branch[0]", false},
		{"$.ledgers[0][22][315].foo", false},
		{"$['store']['book']", false},
		{"$['store'][0]['book']", false},
		{"$.aaa.*", false},
		{"$.aaa[?(@.bbb == 'xxx')]", true},
		{"$.aaa[0,1]", true},
		{"$.aaa[0:1]", true},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := NewReferencePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReferencePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			expr, err := NewPath(tt.path)
			if err != nil {
				t.Errorf("NewPath() error = %v", err)
				return
			}
			if got.Expr.String() != expr.String() {
				t.Errorf("got = '%s', want = '%s'", got.Expr.String(), expr.String())
				return
			}
		})
	}
}
