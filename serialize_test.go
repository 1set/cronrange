package cronrange

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

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
		{"nil struct", args{emptyString, emptyString, 0}, "<nil>"},
		{"empty struct", args{emptyString, emptyString, 0}, emptyString},
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
			if strings.Contains(tt.name, "nil") {
				cr = nil
			} else if tt.args.cronExpr == emptyString {
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

func BenchmarkCronRange_String(b *testing.B) {
	cr, _ := New(exprEveryMin, timeZoneBangkok, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cr.String()
	}
}

func TestCronRange_MarshalJSON(t *testing.T) {
	tempStruct := struct {
		CR    *CronRange
		Name  string
		Value int
	}{
		nil,
		"Test",
		1111,
	}
	type args struct {
		cronExpr    string
		timeZone    string
		durationMin uint64
	}
	tests := []struct {
		name  string
		args  args
		wantJ string
	}{
		{"nil struct", args{emptyString, emptyString, 0}, `{"CR":null,"Name":"Test","Value":1111}`},
		{"empty struct", args{emptyString, emptyString, 0}, `{"CR":null,"Name":"Test","Value":1111}`},
		{"5min duration without time zone", args{exprEveryMin, emptyString, 5}, `{"CR":"DR=5; * * * * *","Name":"Test","Value":1111}`},
		{"10min duration with local time zone", args{exprEveryMin, "local", 10}, `{"CR":"DR=10; * * * * *","Name":"Test","Value":1111}`},
		{"10min duration with time zone", args{exprEveryMin, timeZoneBangkok, 10}, `{"CR":"DR=10; TZ=Asia/Bangkok; * * * * *","Name":"Test","Value":1111}`},
		{"every xmas morning in new york city", args{exprEveryXmasMorning, timeZoneNewYork, 240}, `{"CR":"DR=240; TZ=America/New_York; 0 8 25 12 *","Name":"Test","Value":1111}`},
		{"every new year's day in bangkok", args{exprEveryNewYear, timeZoneBangkok, 1440}, `{"CR":"DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *","Name":"Test","Value":1111}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cr *CronRange
			if strings.Contains(tt.name, "nil") {
				cr = crNil
			} else if tt.args.cronExpr == emptyString {
				cr = crEmpty
			} else {
				var err error
				if cr, err = New(tt.args.cronExpr, tt.args.timeZone, tt.args.durationMin); err != nil {
					t.Errorf("New() error = %v", err)
					return
				}
			}

			tempStruct.CR = cr
			got, err := json.Marshal(tempStruct)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}
			gotJ := string(got)
			if gotJ != tt.wantJ {
				t.Errorf("MarshalJSON() got = %v, want %q", gotJ, tt.wantJ)
			}
		})
	}
}

func BenchmarkCronRange_MarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = crEvery1MinBangkok.MarshalJSON()
	}
}
