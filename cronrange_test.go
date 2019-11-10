package cronrange

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cronExpr    string
		timeZone    string
		durationMin uint64
	}
	tests := []struct {
		name    string
		args    args
		wantCr  bool
		wantErr bool
	}{
		{"Empty cronExpr", args{"", "", 5}, false, true},
		{"Invalid cronExpr", args{"h e l l o", "", 5}, false, true},
		{"Incomplete cronExpr", args{"* * * *", "", 5}, false, true},
		{"Nonexistent timezone", args{"* * * * *", "Mars", 5}, false, true},
		{"Zero durationMin", args{"* * * * *", "", 0}, false, true},
		{"Normal without timezone", args{"* * * * *", "", 5}, true, false},
		{"Normal with local timezone", args{"* * * * *", " Local ", 5}, true, false},
		{"Normal with timezone", args{"* * * * *", "Asia/Bangkok", 5}, true, false},
		{"Normal with large duration", args{"* * * * *", "Asia/Bangkok", 5259000}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCr, err := New(tt.args.cronExpr, tt.args.timeZone, tt.args.durationMin)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotCr != nil) != tt.wantCr {
				t.Errorf("New() gotCr = %v, wantCr %v", gotCr, tt.wantCr)
			}
		})
	}
}
