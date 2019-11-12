package cronrange

import (
	"time"
)

// IsWithin checks if the given time falls within any time range represented by the expression.
func (cr *CronRange) IsWithin(t time.Time) (within bool, err error) {
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

	// if no occurrence is found, it gets zero, i.e. time.Time{}
	if rangeStart.Before(searchStart) {
		within = false
		return
	}

	// check if rangeStart <= t <= rangeEnd
	within = (rangeStart.Before(t) && rangeEnd.After(t)) || rangeStart.Equal(t) || rangeEnd.Equal(t)
	return
}
