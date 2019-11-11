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
		{"nil struct", crNil, "<nil>"},
		{"empty struct", crEmpty, emptyString},
		{"use string() instead of sprintf", crEvery1Min, "DR=1; * * * * *"},
		{"use instance instead of pointer", crEvery1Min, "DR=1; * * * * *"},
		{"1min duration without time zone", crEvery1Min, "DR=1; * * * * *"},
		{"5min duration without time zone", crEvery5Min, "DR=5; */5 * * * *"},
		{"10min duration with local time zone", crEvery10MinLocal, "DR=10; */10 * * * *"},
		{"10min duration with time zone", crEvery10MinBangkok, "DR=10; TZ=Asia/Bangkok; */10 * * * *"},
		{"every xmas morning in new york city", crEveryXmasMorningNYC, "DR=240; TZ=America/New_York; 0 8 25 12 *"},
		{"every new year's day in bangkok", crEveryNewYearsDayBangkok, "DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *"},
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

func TestCronRange_MarshalJSON(t *testing.T) {
	tempStruct := tempTestStruct{
		nil,
		"Test",
		1111,
	}
	tests := []struct {
		name  string
		cr    *CronRange
		wantJ string
	}{
		{"nil struct", crNil, `{"CR":null,"Name":"Test","Value":1111}`},
		{"empty struct", crEmpty, `{"CR":null,"Name":"Test","Value":1111}`},
		{"5min duration without time zone", crEvery5Min, `{"CR":"DR=5; */5 * * * *","Name":"Test","Value":1111}`},
		{"10min duration with local time zone", crEvery10MinLocal, `{"CR":"DR=10; */10 * * * *","Name":"Test","Value":1111}`},
		{"10min duration with time zone", crEvery10MinBangkok, `{"CR":"DR=10; TZ=Asia/Bangkok; */10 * * * *","Name":"Test","Value":1111}`},
		{"every xmas morning in new york city", crEveryXmasMorningNYC, `{"CR":"DR=240; TZ=America/New_York; 0 8 25 12 *","Name":"Test","Value":1111}`},
		{"every new year's day in bangkok", crEveryNewYearsDayBangkok, `{"CR":"DR=1440; TZ=Asia/Bangkok; 0 0 1 1 *","Name":"Test","Value":1111}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempStruct.CR = tt.cr
			got, err := json.Marshal(tempStruct)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}
			gotJ := string(got)
			if gotJ != tt.wantJ {
				t.Errorf("MarshalJSON() got = %v, want %v", gotJ, tt.wantJ)
			}
		})
	}
}

func BenchmarkCronRange_MarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = crEvery10MinBangkok.MarshalJSON()
	}
}

var deserializeTestCases = []struct {
	name    string
	inputS  string
	wantS   string
	wantErr bool
}{
	{"empty string", "", "", true},
	{"invalid expression", "hello", "", true},
	{"missing duration", "; * * * * *", "", true},
	{"invalid duration=0", "DR=0;* * * * *", "", true},
	{"invalid duration=-5", "DR=-5;* * * * *", "", true},
	{"invalid timezone=Mars", "DR=5;TZ=Mars;* * * * *", "", true},
	{"normal without timezone", "DR=5;* * * * *", "DR=5; * * * * *", false},
	{"normal with extra whitespaces", "  DR=6 ;  * * * * *  ", "DR=6; * * * * *", false},
	{"normal with empty parts", ";  DR=7;;; ;; ;; ;* * * * *  ", "DR=7; * * * * *", false},
	{"normal with local time zone", "DR=8;TZ=Local;* * * * *", "DR=8; * * * * *", false},
	{"normal with utc time zone", "DR=9;TZ=Etc/UTC;* * * * *", "DR=9; TZ=Etc/UTC; * * * * *", false},
	{"normal with honolulu time zone", "DR=10;TZ=Pacific/Honolulu;* * * * *", "DR=10; TZ=Pacific/Honolulu; * * * * *", false},
}

func TestCronRange_UnmarshalJSON(t *testing.T) {
	jsonPrefix, jsonSuffix := `{"CR":"`, `","Name":"Demo","Value":2222}`
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			jsonFull := jsonPrefix + tt.inputS + jsonSuffix
			var gotS tempTestStruct
			err := json.Unmarshal([]byte(jsonFull), &gotS)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotS.CR != nil && gotS.CR.String() != tt.wantS {
				t.Errorf("UnmarshalJSON() gotCr: %s, want: %s", gotS.CR.String(), tt.wantS)
			}
			if !tt.wantErr && gotS.CR != nil && (gotS.CR.schedule == nil || gotS.CR.duration == 0) {
				t.Errorf("UnmarshalJSON() incomplete gotCr: %v", gotS.CR)
			}
		})
	}
}

func BenchmarkCronRange_UnmarshalJSON(b *testing.B) {
	jsonFull := []byte(`{"CR":"DR=10;TZ=Pacific/Honolulu;* * * * *","Name":"Demo","Value":2222}`)
	var gotS tempTestStruct
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(jsonFull, &gotS)
	}
}

func TestParseString(t *testing.T) {
	for _, tt := range deserializeTestCases {
		t.Run(tt.name, func(t *testing.T) {
			gotCr, err := ParseString(tt.inputS)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotCr != nil && gotCr.String() != tt.wantS {
				t.Errorf("ParseString() gotCr: %s, want: %s", gotCr.String(), tt.wantS)
			}
			if !tt.wantErr && gotCr != nil && (gotCr.schedule == nil || gotCr.duration == 0) {
				t.Errorf("ParseString() incomplete gotCr: %v", gotCr)
			}
		})
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
		{"from zero to zero", fields{zeroTime, zeroTime}, "[0001-01-01T00:00:00Z,0001-01-01T00:00:00Z]"},
		{"first day of 2020 in utc", fields{firstUtcSec2020, firstUtcSec2020.AddDate(0, 0, 1)}, "[2020-01-01T00:00:00Z,2020-01-02T00:00:00Z]"},
		{"first month of 2019 in bangkok", fields{firstBangkokSec2019, firstBangkokSec2019.AddDate(0, 1, 0)}, "[2019-01-01T00:00:00+07:00,2019-02-01T00:00:00+07:00]"},
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
	tr := TimeRange{firstBangkokSec2019, firstBangkokSec2019.AddDate(0, 1, 0)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.String()
	}
}
