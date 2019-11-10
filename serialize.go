package cronrange

import (
	"encoding/json"
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

type expressionWrapper struct {
	Expression string `json:"expr"`
}

func (cr *CronRange) MarshalJSON() ([]byte, error) {
	expr := cr.String()
	if cr == nil || len(expr) == 0 {
		return []byte("null"), nil
	}

	return json.Marshal(&expressionWrapper {
		Expression:expr,
	})
}

func (cr *CronRange) UnmarshalJSON([]byte) error {
	panic("implement me")
}
