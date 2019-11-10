package cronrange

import (
	"reflect"
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

func TestCronRange_NextOccurrences(t *testing.T) {
	type fields struct {
		Expression  string
		TimeZone    string
		DurationMin uint64
	}
	type args struct {
		t     time.Time
		count int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOccurs []TimeRange
		wantErr    bool
	}{
		{"First 5 minutes in every hour", fields{"0 * * * *", "", 5}, args{time.Now(), 5}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &CronRange{
				Expression:  tt.fields.Expression,
				TimeZone:    tt.fields.TimeZone,
				DurationMin: tt.fields.DurationMin,
			}
			gotOccurs, err := cr.NextOccurrences(tt.args.t, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextOccurrences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOccurs, tt.wantOccurs) {
				t.Errorf("NextOccurrences() gotOccurs = %v, want %v", gotOccurs, tt.wantOccurs)
			}
		})
	}
}