package cronrange

import "time"

var (
	emptyString          = ""
	exprEveryMin         = "* * * * *"
	exprEvery5Min        = "*/5 * * * *"
	exprEvery10Min       = "*/10 * * * *"
	exprEveryXmasMorning = "0 8 25 12 *"
	exprEveryNewYear     = "0 0 1 1 *"
	timeZoneBangkok      = "Asia/Bangkok"
	timeZoneNewYork      = "America/New_York"
	timeZoneHonolulu     = "Pacific/Honolulu"
	timeZoneUTC          = "Etc/UTC"
)

var (
	locationBangkok, _  = time.LoadLocation(timeZoneBangkok)
	locationHonolulu, _ = time.LoadLocation(timeZoneHonolulu)
	locationUTC, _      = time.LoadLocation(timeZoneUTC)
)

var (
	zeroTime          = time.Time{}
	firstLocalSec2012 = time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local)
	firstLocalSec2016 = time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local)
	firstLocalSec2019 = time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
	lastLocalSec2018  = firstLocalSec2019.Add(-1 * time.Second)
)

var (
	crNil                          *CronRange
	crEmpty                        = &CronRange{}
	crEvery1Min, _                 = New(exprEveryMin, emptyString, 1)
	crEvery5Min, _                 = New(exprEvery5Min, emptyString, 5)
	crEvery10MinLocal, _           = New(exprEvery10Min, emptyString, 10)
	crEvery10MinBangkok, _         = New(exprEvery10Min, timeZoneBangkok, 10)
	crEveryXmasMorningNYC, _       = New(exprEveryXmasMorning, timeZoneNewYork, 240)
	crEveryNewYearsDayBangkok, _   = New(exprEveryNewYear, timeZoneBangkok, 1440)
	crFirstDayEachMonth, _         = New("0 0 1 * *", "", 1440)
	crSecondDayEachMonthBangkok, _ = New("0 0 2 * *", timeZoneBangkok, 1440)
	crThirdDayEachMonthHonolulu, _ = New("0 0 3 * *", timeZoneHonolulu, 1440)
	crFirstHourFeb29, _            = New("0 0 29 2 *", "", 60)
	crFirstHourFeb28OrSun, _       = New("0 0 28 2 0", "", 60)
)

type tempTestStruct struct {
	CR    *CronRange
	Name  string
	Value int
}

func getLocalTime(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
}

func getTime(location *time.Location, year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
}

func isTimeRangeSliceEqual(a, b []TimeRange) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if !(a[i].Start.Equal(b[i].Start) && a[i].End.Equal(b[i].End)) {
			return false
		}
	}
	return true
}
