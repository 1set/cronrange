/*
Package cronrange parses Cron-style time range expressions.

In a nutshell, CronRange expression is a combination of Cron expression and time duration to represent periodic time ranges.
And it made easier to figure out if the moment falls within the any time ranges with the IsWithin() method, and what's the next
occurrence with the NextOccurrences() method.

For example, every New Year's Day in Tokyo can be written as:

	DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *

It consists of three parts separated by a semicolon:

    - `DR=1440` stands for duration in minutes, 60 \* 24 = 1440 min;
    - `TZ=Asia/Tokyo` is optional and for time zone using name in IANA Time Zone database (https://www.iana.org/time-zones);
    - `0 0 1 1 *` is a cron expression representing the beginning moment of the time range.

*/
package cronrange
