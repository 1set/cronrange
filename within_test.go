package cronrange

import (
	"testing"
	"time"
)

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
			if tt.crExpr == "nil" {
				cr = nil
			} else if tt.crExpr == "empty" {
				cr = &CronRange{}
			} else {
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
