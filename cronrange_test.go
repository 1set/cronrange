package cronrange

import (
	"fmt"
	"strings"
	"testing"
)

var (
	emptyString          = ""
	exprEveryMin         = "* * * * *"
	exprEveryXmasMorning = "0 8 25 12 *"
	exprEveryNewYear     = "0 0 1 1 *"
	timeZoneBangkok      = "Asia/Bangkok"
	timeZoneNewYork      = "America/New_York"
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

func TestCronRange_String(t *testing.T) {
	type args struct {
		cronExpr    string
		timeZone    string
		durationMin uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil struct", args{emptyString, emptyString, 0}, emptyString},
		{"use string() instead of sprintf", args{exprEveryMin, emptyString, 1}, "DR=1; * * * * *"},
		{"use instance instead of pointer", args{exprEveryMin, emptyString, 1}, "DR=1; * * * * *"},
		{"1min duration without time zone", args{exprEveryMin, emptyString, 1}, "DR=1; * * * * *"},
		{"5min duration without time zone", args{exprEveryMin, emptyString, 5}, "DR=5; * * * * *"},
		{"10min duration with local time zone", args{exprEveryMin, "local", 10}, "DR=10; * * * * *"},
		{"10min duration with time zone", args{exprEveryMin, timeZoneBangkok, 10}, "DR=10; TZ=Asia/Bangkok; * * * * *"},
		{"every xmas morning in new york city", args{exprEveryXmasMorning, timeZoneNewYork, 240}, "DR=240; TZ=America/New_York; 0 8 25 12 *"},
		{"every new year's day in bangkok", args{exprEveryNewYear, timeZoneBangkok, 1440}, "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cr *CronRange
			if tt.args.cronExpr == emptyString {
				cr = &CronRange{}
			} else {
				var err error
				if cr, err = New(tt.args.cronExpr, tt.args.timeZone, tt.args.durationMin); err != nil {
					t.Errorf("New() error = %v", err)
					return
				}
			}
			var got string
			if strings.Contains(tt.name, "string()") {
				got = cr.String()
			} else if strings.Contains(tt.name, "instance") {
				got = fmt.Sprint(*cr)
			} else {
				got = fmt.Sprint(cr)
			}
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkString(b *testing.B) {
	cr, _ := New(exprEveryMin, timeZoneBangkok, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cr.String()
	}
}
