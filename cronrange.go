package cronrange

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// CronRange consists of cron expression along with time zone and duration info.
type CronRange struct {
	cronExpression string
	timeZone       string
	duration       time.Duration
	schedule       cron.Schedule
}

func (cr CronRange) String() string {
	sb := strings.Builder{}
	if cr.duration > 0 {
		sb.WriteString(fmt.Sprintf("DR=%d; ", cr.duration/time.Minute))
	}
	if len(cr.timeZone) > 0 {
		sb.WriteString(fmt.Sprintf("TZ=%s; ", cr.timeZone))
	}
	sb.WriteString(cr.cronExpression)
	return sb.String()
}

// TimeRange represents a time range between starting time and ending time.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("[%v, %v]", tr.Start, tr.End)
}

// New returns a CronRange instance with given config. `timeZone` can be empty for local time zone.
func New(cronExpr, timeZone string, durationMin uint64) (cr *CronRange, err error) {
	// Precondition checks
	if durationMin <= 0 {
		err = errZeroDuration
		return
	}

	// Clean up string parameters
	cronExpr, timeZone = strings.TrimSpace(cronExpr), strings.TrimSpace(timeZone)

	// Append time zone into cron spec if necessary
	cronSpec := cronExpr
	if strings.ToLower(timeZone) == "local" {
		timeZone = ""
	} else if len(timeZone) > 0 {
		cronSpec = fmt.Sprintf("CRON_TZ=%s %s", timeZone, cronExpr)
	}

	// Validate & retrieve crontab schedule
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
