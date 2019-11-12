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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr, err := ParseString(tt.crExpr)
			if err != nil {
				t.Errorf("IsWithin() invalid crExpr: %q, error: %v", cr, err)
				return
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
