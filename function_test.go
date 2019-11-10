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
	crFirstHourFeb29, _ := New("0 0 29 2 *", "", 60)
	crFirstHourFeb28OrSun, _ := New("0 0 28 2 0", "", 60)
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
		{"first hour of feb 29 since 2012",
			crFirstHourFeb29,
			args{time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local), 3},
			[]TimeRange{
				{time.Date(2012, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2012, 2, 29, 1, 0, 0, 0, time.Local)},
				{time.Date(2016, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2016, 2, 29, 1, 0, 0, 0, time.Local)},
				{time.Date(2020, 2, 29, 0, 0, 0, 0, time.Local), time.Date(2020, 2, 29, 1, 0, 0, 0, time.Local)},
			},
			false,
		},
		{"first hour of feb 28 or sunday since 2016",
			crFirstHourFeb28OrSun,
			args{time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local), 6},
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
