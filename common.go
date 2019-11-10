package cronrange

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronZeroTime    = time.Time{}
	cronParseOption = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	cronParser      = cron.NewParser(cronParseOption)
)

var (
	errZeroDuration = errors.New("duration should be positive")
)
