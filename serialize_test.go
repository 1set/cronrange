package cronrange

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCronRange_String(t *testing.T) {
	tests := []struct {
		name string
		cr   *CronRange
		want string
	}{
		{"Nil struct", crNil, "<nil>"},
		{"Empty struct", crEmpty, emptyString},
		{"Use string() instead of sprintf", crEvery1Min, "DR=1; * * * * *"},
		{"Use instance instead of pointer", crEvery1Min, "DR=1; * * * * *"},
		{"1min duration without time zone", crEvery1Min, "DR=1; * * * * *"},
		{"5min duration without time zone", crEvery5Min, "DR=5; */5 * * * *"},
		{"10min duration with local time zone", crEvery10MinLocal, "DR=10; */10 * * * *"},
		{"10min duration with time zone", crEvery10MinBangkok, "DR=10; TZ=Asia/Bangkok; */10 * * * *"},
		{"Every Xmas morning in NYC", crEveryXmasMorningNYC, "DR=240; TZ=America/New_York; 0 8 25 12 *"},
		{"Every New Year's Day in Bangkok", crEveryNewYearsDayBangkok, "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *"},
		{"Every New Year's Day in Tokyo", crEveryNewYearsDayTokyo, "DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := tt.cr
			var got string
			if strings.Contains(tt.name, "string()") {
				got = cr.String()
			} else if strings.Contains(tt.name, "instance") {
				got = fmt.Sprint(*cr)
			} else {
				got = fmt.Sprint(cr)
			}
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkCronRange_String(b *testing.B) {
	cr, _ := New(exprEveryMin, timeZoneBangkok, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cr.String()
	}
}

var deserializeTestCases = []struct {
	name      string
	inputS    string
	wantS     string
	wantErr   bool
	inputOpts []ParseOpt
}{
	{"Empty string", emptyString, emptyString, true, nil},
	{"Invalid expression", "hello", emptyString, true, nil},
	{"Missing duration", "; * * * * *", emptyString, true, nil},
	{"Missing duration default takes priority", "; * * * * *", "DR=1; * * * * *", false, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Missing duration default takes priority fails if invalid", "; * * * * *", emptyString, true, []ParseOpt{DefaultDuration(0 * time.Minute)}},
	{"Invalid duration=0", "DR=0;* * * * *", emptyString, true, nil},
	{"Invalid duration=0 default duration doesn't impact invalid durations", "DR=0;* * * * *", emptyString, true, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Invalid duration=-5", "DR=-5;* * * * *", emptyString, true, nil},
	{"Invalid duration=-5 default doesn't override", "DR=-5;* * * * *", emptyString, true, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Invalid with Mars time zone", "DR=5;TZ=Mars;* * * * *", emptyString, true, nil},
	{"Invalid with unknown part", "DR=10; TZ=Pacific/Honolulu; SET=1; * * * * *", emptyString, true, nil},
	{"Invalid with lower case", "dr=5;* * * * *", emptyString, true, nil},
	{"Invalid with lower case default duration overrides", "dr=5;* * * * *", "DR=1; * * * * *", true, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Invalid with wrong order", "* * * * *; DR=5;", emptyString, true, nil},
	{"Invalid with wrong order fails even with default duration", "* * * * *; DR=5;", emptyString, true, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Normal without timezone", "DR=5;* * * * *", "DR=5; * * * * *", false, nil},
	{"Normal with extra whitespaces", "  DR=6 ;  * * * * *  ", "DR=6; * * * * *", false, nil},
	{"Normal with double duration", "DR=6;DR=7; * * * * *  ", "DR=7; * * * * *", false, nil},
	{"Normal with empty parts", ";  DR=7;;; ;; ;; ;* * * * *  ", "DR=7; * * * * *", false, nil},
	{"Normal with different order", "TZ=Asia/Tokyo;  DR=1440;  0 0 1 1 *", "DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *", false, nil},
	{"Normal with local time zone", "DR=8;TZ=Local;* * * * *", "DR=8; * * * * *", false, nil},
	{"Normal with UTC time zone", "DR=9;TZ=Etc/UTC;* * * * *", "DR=9; TZ=Etc/UTC; * * * * *", false, nil},
	{"Normal with local time zone default duration works with tz", "TZ=Local;* * * * *", "DR=1; * * * * *", false, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Normal with UTC time zone default duration works with tz", "TZ=Etc/UTC;* * * * *", "DR=1; TZ=Etc/UTC; * * * * *", false, []ParseOpt{DefaultDuration(1 * time.Minute)}},
	{"Normal with Honolulu time zone", "DR=10;TZ=Pacific/Honolulu;* * * * *", "DR=10; TZ=Pacific/Honolulu; * * * * *", false, nil},
	{"Normal with Honolulu time zone in different order", "TZ=Pacific/Honolulu; DR=10; * * * * *", "DR=10; TZ=Pacific/Honolulu; * * * * *", false, nil},
	{"Normal with complicated expression", "DR=5258765;   TZ=Pacific/Honolulu;   4,8,22,27,33,38,47,50 3,11,14-16,19,21,22 */10 1,3,5,6,9-11 1-5", "DR=5258765; TZ=Pacific/Honolulu; 4,8,22,27,33,38,47,50 3,11,14-16,19,21,22 */10 1,3,5,6,9-11 1-5", false, nil},
}

func TestParseString(t *testing.T) {
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			gotCr, err := ParseString(tt.inputS, tt.inputOpts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (gotCr == nil || gotCr.schedule == nil || gotCr.duration == 0) {
				t.Errorf("ParseString() incomplete gotCr: %v", gotCr)
				return
			}
			if !tt.wantErr && gotCr.String() != tt.wantS {
				t.Errorf("ParseString() gotCr: %s, want: %s", gotCr.String(), tt.wantS)
			}
		})
	}
}

func BenchmarkParseString(b *testing.B) {
	rs := "DR=10;TZ=Pacific/Honolulu;;* * * * *"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseString(rs)
	}
}

func TestCronRange_MarshalJSON(t *testing.T) {
	tempStructWithPointer := tempTestWithPointer{
		nil,
		"Test",
		1111,
	}
	tempStructWithInstance := tempTestWithInstance{
		CronRange{},
		"Test",
		1111,
	}
	tests := []struct {
		name  string
		cr    *CronRange
		wantJ string
	}{
		{"Nil struct", crNil, `{"CR":null,"Name":"Test","Value":1111}`},
		{"Empty struct", crEmpty, `{"CR":null,"Name":"Test","Value":1111}`},
		{"5min duration without time zone", crEvery5Min, `{"CR":"DR=5; */5 * * * *","Name":"Test","Value":1111}`},
		{"10min duration with local time zone", crEvery10MinLocal, `{"CR":"DR=10; */10 * * * *","Name":"Test","Value":1111}`},
		{"10min duration with time zone", crEvery10MinBangkok, `{"CR":"DR=10; TZ=Asia/Bangkok; */10 * * * *","Name":"Test","Value":1111}`},
		{"Every Xmas morning in NYC", crEveryXmasMorningNYC, `{"CR":"DR=240; TZ=America/New_York; 0 8 25 12 *","Name":"Test","Value":1111}`},
		{"Every New Year's Day in Bangkok", crEveryNewYearsDayBangkok, `{"CR":"DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *","Name":"Test","Value":1111}`},
		{"Every New Year's Day in Tokyo", crEveryNewYearsDayTokyo, `{"CR":"DR=1440; TZ=Asia/Tokyo; 0 0 1 1 *","Name":"Test","Value":1111}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempStructWithPointer.CR = tt.cr
			got, err := json.Marshal(tempStructWithPointer)
			if err != nil {
				t.Errorf("Marshal() with pointer error = %v", err)
				return
			}
			gotJ := string(got)
			if gotJ != tt.wantJ {
				t.Errorf("MarshalJSON() with pointer got = %v, want %v", gotJ, tt.wantJ)
				return
			}

			if tt.cr != nil {
				tempStructWithInstance.CR = *tt.cr
				got, err := json.Marshal(tempStructWithInstance)
				if err != nil {
					t.Errorf("Marshal() with instance error = %v", err)
					return
				}
				gotJ := string(got)
				if gotJ != tt.wantJ {
					t.Errorf("MarshalJSON() with instance got = %v, want %v", gotJ, tt.wantJ)
				}
			}
		})
	}
}

func BenchmarkCronRange_MarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = crEvery10MinBangkok.MarshalJSON()
	}
}

func TestCronRange_UnmarshalJSON_Normal(t *testing.T) {
	jsonPrefix, jsonSuffix := `{"CR":"`, `","Name":"Demo","Value":2222}`
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			jsonFull := jsonPrefix + tt.inputS + jsonSuffix
			var gotSP tempTestWithPointer
			err := json.Unmarshal([]byte(jsonFull), &gotSP)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() with pointer error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (gotSP.CR == nil || gotSP.CR.schedule == nil || gotSP.CR.duration == 0) {
				t.Errorf("UnmarshalJSON() with pointer incomplete gotCr: %v", gotSP.CR)
				return
			}
			if !tt.wantErr && gotSP.CR.String() != tt.wantS {
				t.Errorf("UnmarshalJSON() with pointer gotCr: %s, want: %s", gotSP.CR.String(), tt.wantS)
				return
			}

			var gotSI tempTestWithInstance
			err = json.Unmarshal([]byte(jsonFull), &gotSI)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() with instance error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (gotSI.CR.schedule == nil || gotSI.CR.duration == 0) {
				t.Errorf("UnmarshalJSON() with instance incomplete gotCr: %v", gotSI.CR)
				return
			}
			if !tt.wantErr && gotSI.CR.String() != tt.wantS {
				t.Errorf("UnmarshalJSON() with instance gotCr: %s, want: %s", gotSI.CR.String(), tt.wantS)
				return
			}
		})
	}
}

func TestCronRange_UnmarshalJSON_Broken(t *testing.T) {
	jsonPrefix, jsonSuffix := `{"CR":"`, `","Name":"Demo","Value":2222}`
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			jsonBrokens := []string{
				jsonPrefix[0:len(jsonPrefix)-1] + tt.inputS + jsonSuffix[1:len(jsonSuffix)-1],
				jsonPrefix[0:len(jsonPrefix)-1] + tt.inputS + jsonSuffix,
				jsonPrefix + tt.inputS + jsonSuffix[1:len(jsonSuffix)-1],
				jsonSuffix + jsonPrefix,
				jsonPrefix + tt.inputS,
				tt.inputS + jsonSuffix,
				tt.inputS + jsonPrefix,
				jsonSuffix + tt.inputS,
				jsonSuffix + tt.inputS + jsonPrefix,
				tt.inputS + jsonSuffix + jsonPrefix,
				jsonSuffix + jsonPrefix + tt.inputS,
				tt.inputS + jsonPrefix + jsonSuffix,
				jsonPrefix + jsonSuffix + tt.inputS,
			}
			for _, jsonBroken := range jsonBrokens {
				var gotSP tempTestWithPointer
				if err := json.Unmarshal([]byte(jsonBroken), &gotSP); err == nil {
					t.Errorf("UnmarshalJSON() with pointer missing error for broken json: %s", jsonBroken)
					return
				}

				var gotSI tempTestWithInstance
				if err := json.Unmarshal([]byte(jsonBroken), &gotSI); err == nil {
					t.Errorf("UnmarshalJSON() with instance missing error for broken json: %s", jsonBroken)
					return
				}
			}
		})
	}
}

func TestCronRange_UnmarshalJSON_Direct(t *testing.T) {
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			directExprs := []string{
				`"` + tt.inputS + `"`,
				tt.inputS,
				`"` + tt.inputS,
				tt.inputS + `"`,
			}
			for idx, directExpr := range directExprs {
				var gotCr CronRange
				err := gotCr.UnmarshalJSON([]byte(directExpr))
				if idx == 0 {
					if gotCr.String() != tt.wantS {
						t.Errorf("UnmarshalJSON() directly gotCr: %v, want: %s, expr: %q", gotCr, tt.wantS, directExpr)
						return
					}
				} else {
					if err == nil {
						t.Errorf("UnmarshalJSON() got nil err for: %q", directExpr)
					}
				}
			}
		})
	}
}

func BenchmarkCronRange_UnmarshalJSON(b *testing.B) {
	jsonFull := []byte(`{"CR":"DR=10;TZ=Pacific/Honolulu;* * * * *","Name":"Demo","Value":2222}`)
	var gotS tempTestWithPointer
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(jsonFull, &gotS)
	}
}

func TestTimeRange_String(t *testing.T) {
	type fields struct {
		Start time.Time
		End   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"From zero to zero", fields{zeroTime, zeroTime}, "[0001-01-01T00:00:00Z,0001-01-01T00:00:00Z]"},
		{"First day of 2020 in UTC", fields{firstSec2020Utc, firstSec2020Utc.AddDate(0, 0, 1)}, "[2020-01-01T00:00:00Z,2020-01-02T00:00:00Z]"},
		{"First month of 2019 in Bangkok", fields{firstSec2019Bangkok, firstSec2019Bangkok.AddDate(0, 1, 0)}, "[2019-01-01T00:00:00+07:00,2019-02-01T00:00:00+07:00]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TimeRange{
				Start: tt.fields.Start,
				End:   tt.fields.End,
			}
			if got := tr.String(); got != tt.want {
				t.Errorf("String() = %v, want = %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTimeRange_String(b *testing.B) {
	tr := TimeRange{firstSec2019Bangkok, firstSec2019Bangkok.AddDate(0, 1, 0)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.String()
	}
}
