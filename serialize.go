package cronrange

import (
	"fmt"
	"strings"
	"time"
)

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

func (tr TimeRange) String() string {
	return fmt.Sprintf("[%v, %v]", tr.Start, tr.End)
}
