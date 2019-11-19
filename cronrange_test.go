package cronrange

import (
	"testing"
	"time"
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
		{"Normal with 5 min in Bangkok", args{exprEveryMin, timeZoneBangkok, 5}, true, false},
		{"Normal with 1 day in Tokyo", args{exprEveryNewYear, timeZoneTokyo, 1440}, true, false},
		{"Normal with large duration", args{exprEveryMin, timeZoneBangkok, 5259000}, true, false},
		{"Normal with complicated cron expression", args{exprVeryComplicated, timeZoneHonolulu, 5258765}, true, false},
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

func TestCronRange_Duration(t *testing.T) {
	tests := []struct {
		name    string
		cr      *CronRange
		wantD   time.Duration
		wantErr bool
	}{
		{"Nil struct", crNil, time.Duration(0), true},
		{"Empty struct", crEmpty, time.Duration(0), true},
		{"1min duration without time zone", crEvery1Min, time.Duration(1 * time.Minute), false},
		{"5min duration without time zone", crEvery5Min, time.Duration(5 * time.Minute), false},
		{"10min duration with local time zone", crEvery10MinLocal, time.Duration(10 * time.Minute), false},
		{"Every Xmas morning in NYC", crEveryXmasMorningNYC, time.Duration(4 * time.Hour), false},
		{"Every New Year's Day in Tokyo", crEveryNewYearsDayTokyo, time.Duration(24 * time.Hour), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("Duration() panic = %v, wantErr %v", r, tt.wantErr)
				}
			}()

			got := tt.cr.Duration()
			if got != tt.wantD {
				t.Errorf("Duration() = %v, want %v", got, tt.wantD)
			}
		})
	}
}

func BenchmarkCronRange_Duration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = crEveryNewYearsDayBangkok.Duration()
	}
}

func TestCronRange_TimeZone(t *testing.T) {
	tests := []struct {
		name    string
		cr      *CronRange
		wantS   string
		wantErr bool
	}{
		{"Nil struct", crNil, emptyString, true},
		{"Empty struct", crEmpty, emptyString, true},
		{"5min duration without time zone", crEvery5Min, emptyString, false},
		{"10min duration with local time zone", crEvery10MinLocal, emptyString, false},
		{"Every Xmas morning in NYC", crEveryXmasMorningNYC, "America/New_York", false},
		{"Every New Year's Day in Tokyo", crEveryNewYearsDayTokyo, "Asia/Tokyo", false},
		{"Every the 3rd day in Honolulu", crThirdDayEachMonthHonolulu, "Pacific/Honolulu", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("TimeZone() panic = %v, wantErr %v", r, tt.wantErr)
				}
			}()

			got := tt.cr.TimeZone()
			if got != tt.wantS {
				t.Errorf("TimeZone() = %v, want %v", got, tt.wantS)
			}
		})
	}
}

func BenchmarkCronRange_TimeZone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = crEveryNewYearsDayBangkok.TimeZone()
	}
}

func TestCronRange_CronExpression(t *testing.T) {
	tests := []struct {
		name    string
		cr      *CronRange
		wantS   string
		wantErr bool
	}{
		{"Nil struct", crNil, emptyString, true},
		{"Empty struct", crEmpty, emptyString, true},
		{"1min duration without time zone", crEvery1Min, "* * * * *", false},
		{"5min duration without time zone", crEvery5Min, "*/5 * * * *", false},
		{"10min duration with local time zone", crEvery10MinLocal, "*/10 * * * *", false},
		{"Every Xmas morning in NYC", crEveryXmasMorningNYC, "0 8 25 12 *", false},
		{"Every New Year's Day in Tokyo", crEveryNewYearsDayTokyo, "0 0 1 1 *", false},
		{"Every the 3rd day in Honolulu", crThirdDayEachMonthHonolulu, "0 0 3 * *", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.wantErr {
					t.Errorf("CronExpression() panic = %v, wantErr %v", r, tt.wantErr)
				}
			}()

			got := tt.cr.CronExpression()
			if got != tt.wantS {
				t.Errorf("CronExpression() = %v, want %v", got, tt.wantS)
			}
		})
	}
}

func BenchmarkCronRange_CronExpression(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = crEveryNewYearsDayBangkok.CronExpression()
	}
}
