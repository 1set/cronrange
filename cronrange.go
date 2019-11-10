package cronrange

import (
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronZeroTime = time.Time{}
	cronParseOption = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	cronParser      = cron.NewParser(cronParseOption)
)

type CronRange struct {
	Expression  string
	TimeZone    string
	DurationMin uint64
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func New(expression, timeZone string, durationMin uint64) *CronRange {
	return &CronRange{
		Expression:  expression,
		TimeZone:    timeZone,
		DurationMin: durationMin,
	}
}

func (cr *CronRange) NextOccurrences(t time.Time, count int) (occurs []TimeRange, err error) {
	if cr.DurationMin <= 0 {
		err = errors.New("duration should be positive")
		return
	}

	cronSpec := cr.Expression
	if len(cr.TimeZone) > 0 {
		cronSpec = fmt.Sprintf("CRON_TZ=%s %s", cr.TimeZone, cr.Expression)
	}

	var sched cron.Schedule
	if sched, err = cronParser.Parse(cronSpec); err != nil {
		return
	}

	duration := time.Minute * time.Duration(cr.DurationMin)
	for curr, i := t, 0; i < count; i++ {
		next := sched.Next(curr)
		if next == cronZeroTime {
			break
		}
		occur := TimeRange{
			Start: next,
			End:   next.Add(duration),
		}
		fmt.Println(occur)
		occurs = append(occurs, occur)
		curr = next
	}

	return
}

func NextOccur(spec string) (t time.Time, err error) {
	// spec := "CRON_TZ=Asia/Tokyo 8 * * * 6"
	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := specParser.Parse(spec)
	if err == nil {
		t = sched.Next(time.Now())
	}
	return
}
