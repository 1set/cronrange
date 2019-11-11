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
		{"Empty cronExpr", args{emptyString, emptyString, 5}, false, true},
		{"Invalid cronExpr", args{"h e l l o", emptyString, 5}, false, true},
		{"Incomplete cronExpr", args{"* * * *", emptyString, 5}, false, true},
		{"Nonexistent time zone", args{exprEveryMin, "Mars", 5}, false, true},
		{"Zero durationMin", args{exprEveryMin, emptyString, 0}, false, true},
		{"Normal without time zone", args{exprEveryMin, emptyString, 5}, true, false},
		{"Normal with local time zone", args{exprEveryMin, " Local ", 5}, true, false},
		{"Normal with time zone", args{exprEveryMin, timeZoneBangkok, 5}, true, false},
		{"Normal with large duration", args{exprEveryMin, timeZoneBangkok, 5259000}, true, false},
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

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = New(exprEveryMin, timeZoneBangkok, 10)
	}
}
