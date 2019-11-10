package cronrange

import (
	"testing"
	"time"
)

func TestNextOccur(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		wantT   time.Time
		wantErr bool
	}{
		{"Bangkok Monday Morning", "CRON_TZ=Asia/Bangkok 0 7 * * 1", time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC), false},
		{"Tokyo Monday Morning", "CRON_TZ=Asia/Tokyo 0 9 * * 1", time.Date(2019, 11, 11, 0, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := NextOccur(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextOccur() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotT.UTC() != tt.wantT.UTC() {
				t.Errorf("NextOccur() gotT = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
