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
	errIncompleteExpr     = errors.New("expression should contain at least two parts")
	errMissDurationExpr   = errors.New("duration is missing from the expression")
	errEmptyExpr          = errors.New("got empty expression")
	errJsonNoQuotationFix = errors.New(`json string doesn't start or end with '"'`)
	strDoubleQuotation    = `"`
	strSemicolon          = `;`
	strMarkDuration       = `DR=`
	strMarkTimeZone       = `TZ=`
)

func (cr *CronRange) UnmarshalJSON(b []byte) (err error) {
	raw := string(b)
	if len(raw) == 0 {
		return errEmptyExpr
	}

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
		err = errEmptyExpr
		return
	}

	var (
		cronExpr, timeZone, durStr string
		durMin                     uint64
		parts                      = strings.Split(s, strSemicolon)
		idxExpr                    = len(parts) - 1
	)
	if idxExpr == 0 {
		err = errIncompleteExpr
		return
	}

	for idx, part := range parts {
		part = strings.TrimSpace(part)
		// skip empty part
		if len(part) == 0 {
			continue
		}
		// cron expression must be the last part
		if idx == idxExpr {
			cronExpr = part
		} else if strings.HasPrefix(part, strMarkDuration) {
			durStr = part[len(strMarkDuration):]
			if durMin, err = strconv.ParseUint(durStr, 10, 64); err != nil {
				return
			}
		} else if strings.HasPrefix(part, strMarkTimeZone) {
			timeZone = part[len(strMarkTimeZone):]
		} else {
			err = errors.New(fmt.Sprintf(`json string has unknown part: %q`, part))
			return
		}
	}

	if len(durStr) == 0 {
		err = errMissDurationExpr
		return
	}

	cr, err = New(cronExpr, timeZone, durMin)
	return
}
