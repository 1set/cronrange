package cronrange

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
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

type tempTestStruct struct {
	CR    *CronRange
	Name  string
	Value int
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

func TestCronRange_UnmarshalJSON(t *testing.T) {
	tempStruct := tempTestStruct{
		nil,
		"Demo",
		2222,
	}
	tempStruct.CR = crEvery10MinBangkok
	gotJ, err := json.Marshal(tempStruct)
	fmt.Printf("J: %s, %v\n", gotJ, err)

	var gotS tempTestStruct
	err = json.Unmarshal(gotJ, &gotS)
	fmt.Println(gotS, err)
}
