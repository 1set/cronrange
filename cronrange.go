package cronrange

import (
	"time"

	"github.com/robfig/cron/v3"
)

func NextOccur(spec string) (t time.Time, err error) {
	// spec := "CRON_TZ=Asia/Tokyo 8 * * * 6"
	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := specParser.Parse(spec)
	if err == nil {
		t = sched.Next(time.Now())
	}
	return
}
