package cronrange

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	// TODO: better format with TZ and UTC
	return fmt.Sprintf("[%v, %v]", tr.Start, tr.End)
}

func (cr *CronRange) MarshalJSON() ([]byte, error) {
	expr := cr.String()
	if cr == nil || len(expr) == 0 {
		return []byte("null"), nil
	}

	return json.Marshal(expr)
}

var (
	errEmptyString        = errors.New("got empty string")
	errJsonNoQuotationFix = errors.New(`json string doesn't start or end with '"'`)
	strDoubleQuotation    = `"`
	strSemicolon          = `;`
	strMarkDuration       = `DR=`
	strMarkTimeZone       = `TZ=`
)

func (cr *CronRange) UnmarshalJSON(b []byte) (err error) {
	if b == nil || len(b) == 0 {
		return errEmptyString
	}
	raw := string(b)
	if !(strings.HasPrefix(raw, strDoubleQuotation) && strings.HasSuffix(raw, strDoubleQuotation)) {
		return errJsonNoQuotationFix
	}

	var newCr *CronRange
	if newCr, err = ParseString(raw[1 : len(raw)-1]); err == nil {
		*cr = *newCr
	}
	return
}

func ParseString(s string) (cr *CronRange, err error) {
	if len(s) == 0 {
		err = errEmptyString
		return
	}

	parts := strings.Split(s, strSemicolon)
	var cronExpr, timeZone string
	var durationMin uint64
	idxExpr := len(parts) - 1

	for idx, part := range parts {
		part = strings.TrimSpace(part)

		// cron expression must be the last part
		if idx == idxExpr {
			cronExpr = part
		} else if strings.HasPrefix(part, strMarkDuration) {
			durationStr := part[len(strMarkTimeZone):]
			if durationMin, err = strconv.ParseUint(durationStr, 10, 64); err != nil {
				return
			}
		} else if strings.HasPrefix(part, strMarkTimeZone) {
			timeZone = part[len(strMarkDuration):]
		} else {
			err = errors.New(`json string has unknown part: ` + part)
			return
		}
	}

	cr, err = New(cronExpr, timeZone, durationMin)
	return
}
