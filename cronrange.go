package cronrange

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronZeroTime    = time.Time{}
	cronParseOption = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	cronParser      = cron.NewParser(cronParseOption)
)

type CronRange struct {
	cronExpression string
	timeZone       string
	duration       time.Duration
	schedule       cron.Schedule
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func New(cronExpr, timeZone string, durationMin uint64) (cr *CronRange, err error) {
	if durationMin <= 0 {
		err = errors.New("duration should be positive")
		return
	}

	cronSpec := cronExpr
	if strings.ToLower(strings.TrimSpace(timeZone)) == "local" {
		timeZone = ""
	} else if len(timeZone) > 0 {
		cronSpec = fmt.Sprintf("CRON_TZ=%s %s", timeZone, cronExpr)
	}

	var schedule cron.Schedule
	if schedule, err = cronParser.Parse(cronSpec); err != nil {
		return
	}

	cr = &CronRange{
		cronExpression: cronExpr,
		timeZone:       timeZone,
		duration:       time.Minute * time.Duration(durationMin),
		schedule:       schedule,
	}
	return
}
