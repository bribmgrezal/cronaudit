package schedule

import (
	"testing"
)

func TestValidate_ValidExpressions(t *testing.T) {
	cases := []string{
		"* * * * *",
		"0 * * * *",
		"*/5 * * * *",
		"0 9-17 * * 1-5",
		"0,30 6 * * *",
		"15 14 1 * *",
		"0 0 * * 0",
		"*/15 */2 * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			if errs := Validate(expr); len(errs) != 0 {
				t.Errorf("expected no errors, got %v", errs)
			}
		})
	}
}

func TestValidate_WrongFieldCount(t *testing.T) {
	if errs := Validate("* * * *"); len(errs) == 0 {
		t.Error("expected error for 4 fields")
	}
	if errs := Validate("* * * * * *"); len(errs) == 0 {
		t.Error("expected error for 6 fields")
	}
}

func TestValidate_MinuteOutOfRange(t *testing.T) {
	errs := Validate("60 * * * *")
	if len(errs) == 0 {
		t.Error("expected error for minute=60")
	}
}

func TestValidate_HourOutOfRange(t *testing.T) {
	errs := Validate("0 24 * * *")
	if len(errs) == 0 {
		t.Error("expected error for hour=24")
	}
}

func TestValidate_InvalidStep(t *testing.T) {
	cases := []string{
		"*/0 * * * *",
		"*/abc * * * *",
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			if errs := Validate(expr); len(errs) == 0 {
				t.Errorf("expected error for %q", expr)
			}
		})
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cases := []string{
		"0-60 * * * *",  // minute hi out of range
		"10-5 * * * *",  // lo > hi
		"0 25-30 * * *", // hour hi out of range
	}
	for _, expr := range cases {
		t.Run(expr, func(t *testing.T) {
			if errs := Validate(expr); len(errs) == 0 {
				t.Errorf("expected error for %q", expr)
			}
		})
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	// both minute and hour are invalid
	errs := Validate("99 99 * * *")
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(errs))
	}
}

func TestValidate_ListValues(t *testing.T) {
	if errs := Validate("0,15,30,45 * * * *"); len(errs) != 0 {
		t.Errorf("valid list should not produce errors: %v", errs)
	}
	if errs := Validate("0,61 * * * *"); len(errs) == 0 {
		t.Error("expected error for out-of-range list value")
	}
}
