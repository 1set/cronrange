package cronrange

import (
	"testing"
	"time"
)

func TestCronRange_NextOccurrences(t *testing.T) {
	type args struct {
		t     time.Time
		count int
	}
	tests := []struct {
		name       string
		cr         *CronRange
		args       args
		wantOccurs []TimeRange
		wantErr    bool
	}{
		{"Nil struct",
			crNil,
			args{firstSec2019Local, 1},
			nil,
			true,
		},
		{"Empty struct",
			crEmpty,
			args{firstSec2019Local, 1},
			nil,
			true,
		},
		{"Zero count",
			crFirstDayEachMonth,
			args{firstSec2019Local, 0},
			nil,
			true,
		},
		{"Negative count",
			crFirstDayEachMonth,
			args{firstSec2019Local, -5},
			nil,
			true,
		},
		{"Every New Year's Day in Tokyo in since 2018",
			crEveryNewYearsDayTokyo,
			args{firstSec2018Tokyo.Add(-1 * time.Second), 3},
			[]TimeRange{
				{parseTime(locationTokyo, "2018-01-01 00:00:00"), parseTime(locationTokyo, "2018-01-02 00:00:00")},
				{parseTime(locationTokyo, "2019-01-01 00:00:00"), parseTime(locationTokyo, "2019-01-02 00:00:00")},
				{parseTime(locationTokyo, "2020-01-01 00:00:00"), parseTime(locationTokyo, "2020-01-02 00:00:00")},
			},
			false,
		},
		{"First day of january in 2019",
			crFirstDayEachMonth,
			args{lastSec2018Local, 1},
			[]TimeRange{
				{firstSec2019Local, parseLocalTime("2019-01-02 00:00:00")},
			},
			false,
		},
		{"First day of january in 2019 in Honolulu",
			crFirstDayEachMonth,
			args{parseTime(locationHonolulu, "2019-01-01 00:00:00").Add(-1 * time.Second), 1},
			[]TimeRange{
				{parseTime(locationHonolulu, "2019-01-01 00:00:00"), parseTime(locationHonolulu, "2019-01-02 00:00:00")},
			},
			false,
		},
		{"Second day of january in 2019 in Bangkok (UTC view)",
			crSecondDayEachMonthBangkok,
			args{firstSec2019Local, 1},
			[]TimeRange{
				{parseTime(locationUTC, "2019-01-01 17:00:00"), parseTime(locationUTC, "2019-01-02 17:00:00")},
			},
			false,
		},
		{"Third day of january in 2019 in Honolulu (Bangkok view)",
			crThirdDayEachMonthHonolulu,
			args{parseTime(locationBangkok, "2019-01-01 00:00:00"), 1},
			[]TimeRange{
				{parseTime(locationBangkok, "2019-01-03 17:00:00"), parseTime(locationBangkok, "2019-01-04 17:00:00")},
			},
			false,
		},
		{"First day of first three months in 2019",
			crFirstDayEachMonth,
			args{lastSec2018Local, 3},
			[]TimeRange{
				{firstSec2019Local, parseLocalTime("2019-01-02 00:00:00")},
				{parseLocalTime("2019-02-01 00:00:00"), parseLocalTime("2019-02-02 00:00:00")},
				{parseLocalTime("2019-03-01 00:00:00"), parseLocalTime("2019-03-02 00:00:00")},
			},
			false,
		},
		{"First hour of feb 29 since 2012",
			crFirstHourFeb29,
			args{firstSec2012Local, 3},
			[]TimeRange{
				{parseLocalTime("2012-02-29 00:00:00"), parseLocalTime("2012-02-29 01:00:00")},
				{parseLocalTime("2016-02-29 00:00:00"), parseLocalTime("2016-02-29 01:00:00")},
				{parseLocalTime("2020-02-29 00:00:00"), parseLocalTime("2020-02-29 01:00:00")},
			},
			false,
		},
		{"First hour of feb 28 or sunday since 2016",
			crFirstHourFeb28OrSun,
			args{firstSec2016Local, 6},
			[]TimeRange{
				{parseLocalTime("2016-02-07 00:00:00"), parseLocalTime("2016-02-07 01:00:00")},
				{parseLocalTime("2016-02-14 00:00:00"), parseLocalTime("2016-02-14 01:00:00")},
				{parseLocalTime("2016-02-21 00:00:00"), parseLocalTime("2016-02-21 01:00:00")},
				{parseLocalTime("2016-02-28 00:00:00"), parseLocalTime("2016-02-28 01:00:00")},
				{parseLocalTime("2017-02-05 00:00:00"), parseLocalTime("2017-02-05 01:00:00")},
				{parseLocalTime("2017-02-12 00:00:00"), parseLocalTime("2017-02-12 01:00:00")},
			},
			false,
		},
		{"First days of jan with overlap in 2012",
			crEveryDayWithOverlap,
			args{firstSec2012Local, 5},
			[]TimeRange{
				{parseLocalTime("2012-01-02 00:00:00"), parseLocalTime("2012-01-04 00:00:00")},
				{parseLocalTime("2012-01-03 00:00:00"), parseLocalTime("2012-01-05 00:00:00")},
				{parseLocalTime("2012-01-04 00:00:00"), parseLocalTime("2012-01-06 00:00:00")},
				{parseLocalTime("2012-01-05 00:00:00"), parseLocalTime("2012-01-07 00:00:00")},
				{parseLocalTime("2012-01-06 00:00:00"), parseLocalTime("2012-01-08 00:00:00")},
			},
			false,
		},
		{"Very complicated time periods since 2017",
			crVeryComplicated,
			args{firstSec2017Honolulu, 5},
			[]TimeRange{
				{parseTime(locationHonolulu, "2017-01-01 03:04:00"), parseTime(locationHonolulu, "2017-01-02 01:41:00")},
				{parseTime(locationHonolulu, "2017-01-01 03:08:00"), parseTime(locationHonolulu, "2017-01-02 01:45:00")},
				{parseTime(locationHonolulu, "2017-01-01 03:22:00"), parseTime(locationHonolulu, "2017-01-02 01:59:00")},
				{parseTime(locationHonolulu, "2017-01-01 03:27:00"), parseTime(locationHonolulu, "2017-01-02 02:04:00")},
				{parseTime(locationHonolulu, "2017-01-01 03:33:00"), parseTime(locationHonolulu, "2017-01-02 02:10:00")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOccurs, err := tt.cr.NextOccurrences(tt.args.t, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextOccurrences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !isTimeRangeSliceEqual(gotOccurs, tt.wantOccurs) {
				t.Errorf("NextOccurrences() gotOccurs = %v, want %v", gotOccurs, tt.wantOccurs)
			}
		})
	}
}

func BenchmarkCronRange_NextOccurrences(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = crEvery10MinBangkok.NextOccurrences(firstSec2019Local, 10)
	}
}

func TestCronRange_IsWithin(t *testing.T) {
	tests := []struct {
		name       string
		crExpr     string
		t          time.Time
		wantWithin bool
		wantErr    bool
	}{
		{"Nil instance", "nil", parseLocalTime("2019-01-01 01:00:30"), false, true},
		{"Empty instance", "empty", parseLocalTime("2019-01-01 01:00:30"), false, true},
		{"Every 3rd minute - in", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 01:00:30"), true, false},
		{"Every 3rd minute - out1", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 01:02:00"), false, false},
		{"Every 3rd minute - out2", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 00:59:59"), false, false},
		{"Every 3rd minute - out3", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 01:01:01"), false, false},
		{"Every 3rd minute - lower", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 01:00:00"), true, false},
		{"Every 3rd minute - upper", "DR=1; */3 * * * *", parseLocalTime("2019-01-01 01:01:00"), true, false},
		{"Every 3rd hour - in1", "DR=60; 0 */3 * * *", parseLocalTime("2019-01-01 00:25:00"), true, false},
		{"Every 3rd hour - in2", "DR=60; 0 */3 * * *", parseLocalTime("2019-01-01 00:55:00"), true, false},
		{"Every 3rd hour - out1", "DR=60; 0 */3 * * *", parseLocalTime("2019-01-01 01:00:01"), false, false},
		{"Every 3rd hour - out2", "DR=60; 0 */3 * * *", parseLocalTime("2019-01-01 02:59:00"), false, false},
		{"Every continuous hour - 1", "DR=60; 0 * * * *", parseLocalTime("2019-01-01 00:00:00"), true, false},
		{"Every continuous hour - 2", "DR=60; 0 * * * *", parseLocalTime("2019-01-01 00:59:59"), true, false},
		{"Every continuous hour - 3", "DR=60; 0 * * * *", parseLocalTime("2019-01-01 01:00:00"), true, false},
		{"Every continuous hour - 4", "DR=60; 0 * * * *", parseLocalTime("2019-01-01 01:00:01"), true, false},
		{"Every continuous hour - 5", "DR=60; 0 * * * *", parseLocalTime("2019-01-01 02:00:00"), true, false},
		{"Every continuous hour - 6", "DR=60; 0 * * * *", parseLocalTime("2020-01-23 04:56:19"), true, false},
		{"Every hour with overlap - 1", "DR=65; 0 * * * *", parseLocalTime("2019-01-01 00:00:00"), true, false},
		{"Every hour with overlap - 2", "DR=65; 0 * * * *", parseLocalTime("2019-01-01 00:59:59"), true, false},
		{"Every hour with overlap - 3", "DR=65; 0 * * * *", parseLocalTime("2019-01-01 01:00:00"), true, false},
		{"Every hour with overlap - 4", "DR=65; 0 * * * *", parseLocalTime("2019-01-01 01:00:01"), true, false},
		{"Every hour with overlap - 5", "DR=65; 0 * * * *", parseLocalTime("2019-01-01 02:00:00"), true, false},
		{"Every hour with overlap - 6", "DR=65; 0 * * * *", parseLocalTime("2020-01-23 04:56:19"), true, false},
		{"Every New Year's Day - in", "DR=1440; 0 0 1 1 *", parseLocalTime("2019-01-01 12:34:56"), true, false},
		{"Every New Year's Day - out1", "DR=1440; 0 0 1 1 *", parseLocalTime("2019-02-01 12:34:56"), false, false},
		{"Every New Year's Day - out2", "DR=1440; 0 0 1 1 *", parseLocalTime("2019-01-02 00:00:01"), false, false},
		{"Every New Year's Day in Bangkok - in1", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationBangkok, "2019-01-01 12:34:56"), true, false},
		{"Every New Year's Day in Bangkok - in2", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationBangkok, "2019-01-02 00:00:00"), true, false},
		{"Every New Year's Day in Bangkok - in3", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationUTC, "2019-01-01 17:00:00"), true, false},
		{"Every New Year's Day in Bangkok - in4", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationUTC, "2019-01-01 00:00:00"), true, false},
		{"Every New Year's Day in Bangkok - out1", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationBangkok, "2019-01-02 00:00:01"), false, false},
		{"Every New Year's Day in Bangkok - out2", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationBangkok, "2019-01-03 00:00:00"), false, false},
		{"Every New Year's Day in Bangkok - out3", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationUTC, "2019-01-01 17:00:01"), false, false},
		{"Every New Year's Day in Bangkok - out4", "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *", parseTime(locationUTC, "2019-01-02 00:00:00"), false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cr *CronRange
			switch {
			case tt.crExpr == "nil":
				cr = nil
			case tt.crExpr == "empty":
				cr = &CronRange{}
			default:
				var err error
				if cr, err = ParseString(tt.crExpr); err != nil {
					t.Errorf("IsWithin() invalid crExpr: %q, error: %v", cr, err)
					return
				}
			}

			gotWithin, err := cr.IsWithin(tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsWithin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWithin != tt.wantWithin {
				t.Errorf("IsWithin() gotWithin = %v, want %v", gotWithin, tt.wantWithin)
			}
		})
	}
}

func BenchmarkCronRange_IsWithin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = crEvery10MinBangkok.IsWithin(firstSec2019Bangkok)
	}
}
