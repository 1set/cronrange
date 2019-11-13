package cronrange

import (
	"errors"
	"time"
)

var (
	errZeroOrNegCount   = errors.New("count should be positive")
	errNilCronRange     = errors.New("nil CronRange instance")
	errInvalidCronRange = errors.New("invalid CronRange instance")
)

// NextOccurrences returns the next occurrence time ranges, later than the given time.
func (cr *CronRange) NextOccurrences(t time.Time, count int) (occurs []TimeRange, err error) {
	// Precondition checks
	switch {
	case count <= 0:
		err = errZeroOrNegCount
	case cr == nil:
		err = errNilCronRange
	case cr.schedule == nil, cr.duration <= 0:
		err = errInvalidCronRange
	}
	if err != nil {
		return
	}

	for curr, i := t, 0; i < count; i++ {
		// if no occurrence is found within next five years, it returns zero time, i.e. time.Time{}
		next := cr.schedule.Next(curr)
		if next.Before(curr) {
			break
		}
		occur := TimeRange{
			Start: next,
			End:   next.Add(cr.duration),
		}
		occurs = append(occurs, occur)
		curr = next
	}

	return
}

// IsWithin checks if the given time falls within any time range represented by the expression.
func (cr *CronRange) IsWithin(t time.Time) (within bool, err error) {
	within = false
	if cr == nil {
		err = errNilCronRange
		return
	} else if cr.schedule == nil || cr.duration < 0 {
		err = errInvalidCronRange
		return
	}

	searchStart := t.Add(-(cr.duration + 1*time.Second))
	rangeStart := cr.schedule.Next(searchStart)
	rangeEnd := rangeStart.Add(cr.duration)

	// if no occurrence is found, it gets zero time, i.e. time.Time{}
	if rangeStart.Before(searchStart) {
		return
	}

	// check if rangeStart <= t <= rangeEnd
	within = (rangeStart.Before(t) && rangeEnd.After(t)) || rangeStart.Equal(t) || rangeEnd.Equal(t)
	return
}
