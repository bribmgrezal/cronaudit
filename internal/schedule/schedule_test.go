package schedule

import (
	"testing"
)

func TestParse_Wildcard(t *testing.T) {
	s, err := Parse("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Minute.Values) != 60 {
		t.Errorf("expected 60 minute values, got %d", len(s.Minute.Values))
	}
	if len(s.Hour.Values) != 24 {
		t.Errorf("expected 24 hour values, got %d", len(s.Hour.Values))
	}
}

func TestParse_Step(t *testing.T) {
	s, err := Parse("*/15 */6 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantMinutes := []int{0, 15, 30, 45}
	for i, v := range s.Minute.Values {
		if v != wantMinutes[i] {
			t.Errorf("minute[%d]: want %d, got %d", i, wantMinutes[i], v)
		}
	}
	wantHours := []int{0, 6, 12, 18}
	for i, v := range s.Hour.Values {
		if v != wantHours[i] {
			t.Errorf("hour[%d]: want %d, got %d", i, wantHours[i], v)
		}
	}
}

func TestParse_Range(t *testing.T) {
	s, err := Parse("0 9-17 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Hour.Values) != 9 {
		t.Errorf("expected 9 hour values, got %d", len(s.Hour.Values))
	}
	if len(s.DayOfWeek.Values) != 5 {
		t.Errorf("expected 5 dow values, got %d", len(s.DayOfWeek.Values))
	}
}

func TestParse_List(t *testing.T) {
	s, err := Parse("0,30 8,12,18 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Minute.Values) != 2 {
		t.Errorf("expected 2 minute values, got %d", len(s.Minute.Values))
	}
	if len(s.Hour.Values) != 3 {
		t.Errorf("expected 3 hour values, got %d", len(s.Hour.Values))
	}
}

func TestParse_InvalidFieldCount(t *testing.T) {
	_, err := Parse("* * * *")
	if err == nil {
		t.Error("expected error for 4-field expression")
	}
}

func TestParse_OutOfRange(t *testing.T) {
	cases := []string{
		"60 * * * *",
		"* 24 * * *",
		"* * 0 * *",
		"* * * 13 *",
	}
	for _, c := range cases {
		_, err := Parse(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestParse_InvalidStep(t *testing.T) {
	_, err := Parse("*/0 * * * *")
	if err == nil {
		t.Error("expected error for step 0")
	}
}
