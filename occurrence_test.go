package cronrange

import (
	"reflect"
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
		{"nil struct",
			crNil,
			args{firstLocalSec2019, 1},
			nil,
			true,
		},
		{"empty struct",
			crEmpty,
			args{firstLocalSec2019, 1},
			nil,
			true,
		},
		{"zero count",
			crFirstDayEachMonth,
			args{firstLocalSec2019, 0},
			nil,
			true,
		},
		{"negative count",
			crFirstDayEachMonth,
			args{firstLocalSec2019, -5},
			nil,
			true,
		},
		{"first day of first month in 2019",
			crFirstDayEachMonth,
			args{lastLocalSec2018, 1},
			[]TimeRange{
				{firstLocalSec2019, time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
			},
			false,
		},
		{"first day of first three months in 2019",
			crFirstDayEachMonth,
			args{lastLocalSec2018, 3},
			[]TimeRange{
				{firstLocalSec2019, time.Date(2019, 1, 2, 0, 0, 0, 0, time.Local)},
				{time.Date(2019, 2, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 2, 2, 0, 0, 0, 0, time.Local)},
				{time.Date(2019, 3, 1, 0, 0, 0, 0, time.Local), time.Date(2019, 3, 2, 0, 0, 0, 0, time.Local)},
			},
			false,
		},
		{"first hour of feb 29 since 2012",
			crFirstHourFeb29,
			args{firstLocalSec2012, 3},
			[]TimeRange{
				{time.Date(2012, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2012, 2, 29, 1, 0, 0, 0, time.Local)},
				{time.Date(2016, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 29, 1, 0, 0, 0, time.Local)},
				{time.Date(2020, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2020, 2, 29, 1, 0, 0, 0, time.Local)},
			},
			false,
		},
		{"first hour of feb 28 or sunday since 2016",
			crFirstHourFeb28OrSun,
			args{firstLocalSec2016, 6},
			[]TimeRange{
				{time.Date(2016, 2, 7, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 7, 1, 0, 0, 0, time.Local)},
				{time.Date(2016, 2, 14, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 14, 1, 0, 0, 0, time.Local)},
				{time.Date(2016, 2, 21, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 21, 1, 0, 0, 0, time.Local)},
				{time.Date(2016, 2, 28, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 28, 1, 0, 0, 0, time.Local)},
				{time.Date(2017, 2, 5, 0, 0, 0, 0, time.Local), time.Date(2017, 2, 5, 1, 0, 0, 0, time.Local)},
				{time.Date(2017, 2, 12, 0, 0, 0, 0, time.Local), time.Date(2017, 2, 12, 1, 0, 0, 0, time.Local)},
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

// TODO: Benchmark for NC
