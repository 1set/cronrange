package cronrange

import (
	"reflect"
	"testing"
	"time"
)

func TestCronRange_NextOccurrences(t *testing.T) {
	var crNil *CronRange
	crEmpty := &CronRange{}
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
		{"nil struct",
			crNil,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), 1},
			nil,
			true,
		},
		{"empty struct",
			crEmpty,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), 1},
			nil,
			true,
		},
		{"zero count",
			crFirstDayEachMonth,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), 0},
			nil,
			true,
		},
		{"negative count",
			crFirstDayEachMonth,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), -5},
			nil,
			true,
		},
		{"first day of first month in 2019",
			crFirstDayEachMonth,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second), 1},
			[]TimeRange{
				{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
			},
			false,
		},
		{"first day of first three months in 2019",
			crFirstDayEachMonth,
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second), 3},
			[]TimeRange{
				{time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
				{time.Date(2019, 2, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 2, 2, 0, 0, 0, 0, time.Local)},
				{time.Date(2019, 3, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 3, 2, 0, 0, 0, 0, time.Local)},
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
