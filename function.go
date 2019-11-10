package cronrange

import (
	"time"
)

func (cr *CronRange) NextOccurrences(t time.Time, count int) (occurs []TimeRange, err error) {
	if count <= 0 {
		err = errZeroOrNegCount
	} else if cr == nil {
		err = errNilCronRange
	} else if cr.schedule == nil || cr.duration < 0 {
		err = errInvalidCronRange
	}

	if err != nil {
		return
	}

	for curr, i := t, 0; i < count; i++ {
		next := cr.schedule.Next(curr)
		if next == cronZeroTime {
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
