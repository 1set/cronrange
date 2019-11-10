package cronrange

import (
	"time"
)

func (cr *CronRange) NextOccurrences(t time.Time, count int) (occurs []TimeRange, err error) {
	if cr == nil {
		err = errNilCronRange
		return
	} else if count <= 0 {
		err = errZeroOrNegCount
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
