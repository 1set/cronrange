package cronrange

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	strSingleWhitespace = ` `
	strDoubleQuotation  = `"`
	strSemicolon        = `;`
	strMarkDuration     = `DR=`
	strMarkTimeZone     = `TZ=`

	errIncompleteExpr     = errors.New("expression should contain at least two parts")
	errMissDurationExpr   = errors.New("duration is missing from the expression")
	errEmptyExpr          = errors.New("expression is empty")
	errJSONNoQuotationFix = errors.New(`json string should start and end with '"'`)
)

// String returns a normalized CronRange expression, which can be consumed by ParseString().
func (cr CronRange) String() string {
	sb := strings.Builder{}
	sb.Grow(36)
	if cr.duration > 0 {
		sb.WriteString(strMarkDuration)
		if cr.version == Version2 {
			sb.WriteString(cr.duration.String())
		} else {
			sb.WriteString(strconv.FormatUint(uint64(cr.duration/time.Minute), 10))
		}
		sb.WriteString(strSemicolon)
		sb.WriteString(strSingleWhitespace)
	}
	if len(cr.timeZone) > 0 {
		sb.WriteString(strMarkTimeZone)
		sb.WriteString(cr.timeZone)
		sb.WriteString(strSemicolon)
		sb.WriteString(strSingleWhitespace)
	}
	sb.WriteString(cr.cronExpression)
	return sb.String()
}

// ParseString attempts to deserialize the given expression or return failure if any parsing errors occur.
func ParseString(s string) (cr *CronRange, err error) {
	cr, err = internalParseString(s, cronParser)
	return
}

func internalParseString(s string, cp cron.Parser) (cr *CronRange, err error) {
	if s == "" {
		err = errEmptyExpr
		return
	}

	var (
		cronExpr, timeZone, durStr string
		durMin                     uint64
		parts                      = strings.Split(s, strSemicolon)
		idxExpr                    = len(parts) - 1
		version                    int
		duration                   time.Duration
	)
	if idxExpr == 0 {
		err = errIncompleteExpr
		return
	}

PL:
	for idx, part := range parts {
		part = strings.TrimSpace(part)
		// skip empty part
		if part == "" {
			continue
		}

		switch {
		case idx == idxExpr:
			// cron expression must be the last part
			cronExpr = part
		case strings.HasPrefix(part, strMarkDuration):
			durStr = part[len(strMarkDuration):]
			if duration, err = time.ParseDuration(durStr); err == nil {
				if duration > 0 {
					version = Version2
				}
			}
			if durMin, err = strconv.ParseUint(durStr, 10, 64); err != nil {
				if version != Version2 {
					break PL
				}
				err = nil
			}
		case strings.HasPrefix(part, strMarkTimeZone):
			timeZone = part[len(strMarkTimeZone):]
		default:
			err = fmt.Errorf(`expression got unknown part: %q`, part)
		}
	}

	if err == nil {
		if len(durStr) > 0 {
			if version == Version2 {
				cr, err = Create(cronExpr, timeZone, duration, cp)
			} else {
				cr, err = internalNew(cronExpr, timeZone, time.Duration(durMin)*time.Minute, cp)
			}
		} else {
			err = errMissDurationExpr
		}
	}
	return
}

func ParseStringWithCronParser(s string, cp cron.Parser) (cr *CronRange, err error) {
	cr, err = internalParseString(s, cp)
	return
}

// MarshalJSON implements the encoding/json.Marshaler interface for serialization of CronRange.
func (cr CronRange) MarshalJSON() ([]byte, error) {
	expr := cr.String()
	if expr == "" {
		return []byte("null"), nil
	}
	return json.Marshal(expr)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface for deserialization of CronRange.
func (cr *CronRange) UnmarshalJSON(b []byte) (err error) {
	// Precondition checks
	raw := string(b)
	if raw == "" {
		return errEmptyExpr
	}
	if !(strings.HasPrefix(raw, strDoubleQuotation) && strings.HasSuffix(raw, strDoubleQuotation) && len(raw) >= 2) {
		return errJSONNoQuotationFix
	}

	// Extract and treat as CronRange expression
	var newCr *CronRange
	if newCr, err = ParseString(raw[1 : len(raw)-1]); err == nil {
		*cr = *newCr
	}
	return
}

// String returns a string representing time range with formatted time values in Internet RFC 3339 format.
func (tr TimeRange) String() string {
	sb := strings.Builder{}
	sb.Grow(54)
	sb.WriteString("[")
	sb.WriteString(tr.Start.Format(time.RFC3339))
	sb.WriteString(",")
	sb.WriteString(tr.End.Format(time.RFC3339))
	sb.WriteString("]")
	return sb.String()
}
