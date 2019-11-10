package cronrange

import (
	"reflect"
	"testing"
	"time"
)

func TestCronRange_NextOccurrences(t *testing.T) {
	crFirstDayEachMonth, _ := New("0 0 1 * *", "", 1440)
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
		{"First day of first month in 2019",
			crFirstDayEachMonth,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second), 1},
			[]TimeRange{
				{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
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
			if !reflect.DeepEqual(gotOccurs, tt.wantOccurs) {
				t.Errorf("NextOccurrences() gotOccurs = %v, want %v", gotOccurs, tt.wantOccurs)
			}
		})
	}
}
