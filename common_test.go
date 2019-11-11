package cronrange

import "time"

var (
	emptyString          = ""
	exprEveryMin         = "* * * * *"
	exprEveryXmasMorning = "0 8 25 12 *"
	exprEveryNewYear     = "0 0 1 1 *"
	timeZoneBangkok      = "Asia/Bangkok"
	timeZoneNewYork      = "America/New_York"
	timeZoneHonolulu     = "Pacific/Honolulu"
	timeZoneUTC          = "Etc/UTC"
)

var (
	firstLocalSec2012 = time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local)
	firstLocalSec2016 = time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local)
	firstLocalSec2019 = time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)
	lastLocalSec2018  = firstLocalSec2019.Add(-1 * time.Second)
)

var (
	crNil                    *CronRange
	crEmpty                  = &CronRange{}
	crEvery1Min, _           = New(exprEveryMin, emptyString, 1)
	crEvery1MinBangkok, _    = New(exprEveryMin, timeZoneBangkok, 10)
	crFirstDayEachMonth, _   = New("0 0 1 * *", "", 1440)
	crFirstHourFeb29, _      = New("0 0 29 2 *", "", 60)
	crFirstHourFeb28OrSun, _ = New("0 0 28 2 0", "", 60)
)

func getLocalTime(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
}

func getTime(location *time.Location, year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
}
