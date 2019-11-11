package cronrange

import (
	"errors"

	"github.com/robfig/cron/v3"
)

var (
	cronParseOption = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	cronParser      = cron.NewParser(cronParseOption)
)

var (
	errZeroDuration     = errors.New("duration should be positive")
	errZeroOrNegCount   = errors.New("count should be positive")
	errNilCronRange     = errors.New("nil CronRange instance")
	errInvalidCronRange = errors.New("invalid CronRange instance")
)
