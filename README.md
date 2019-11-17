# cronrange

[![GoDoc](https://godoc.org/github.com/1set/cronrange?status.svg)](https://godoc.org/github.com/1set/cronrange)
[![License](https://img.shields.io/github/license/1set/cronrange)](https://github.com/1set/cronrange/blob/master/LICENSE)
[![GitHub Action Workflow](https://github.com/1set/cronrange/workflows/build/badge.svg)](https://github.com/1set/cronrange/actions?workflow=build)
[![Go Report Card](https://goreportcard.com/badge/github.com/1set/cronrange)](https://goreportcard.com/report/github.com/1set/cronrange)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ef272059b4044252b0097270b48d5703)](https://www.codacy.com/manual/an9an63/cronrange)
[![Codecov](https://codecov.io/gh/1set/cronrange/branch/master/graph/badge.svg)](https://codecov.io/gh/1set/cronrange)

cronrange is a Go package for _time range expression_ in _cron_ style, or rather:

> CronRange := Cron + TimeRange

for example, every New Year's Day in Tokyo can be written as:

```cron
DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *
```

-   `DR=1440` stands for duration in minutes, 60 \* 24 = 1440 min;
-   `TZ=Asia/Tokyo` is optional and for time zone using name in [IANA Time Zone database](https://www.iana.org/time-zones);
-   `0 0 1 1 *` is a cron expression representing the beginning of the time range.

* * *

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2F1set%2Fcronrange.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2F1set%2Fcronrange?ref=badge_large)
