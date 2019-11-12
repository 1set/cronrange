package cronrange

import (
	"testing"
	"time"
)

func TestCronRange_IsWithin(t *testing.T) {
	tests := []struct {
		name       string
		cr         *CronRange
		t          time.Time
		wantWithin bool
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := tt.cr
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
