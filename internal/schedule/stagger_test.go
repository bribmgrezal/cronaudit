package schedule

import (
	"strings"
	"testing"
)

func TestStagger_WildcardMinute(t *testing.T) {
	res, err := Stagger("* * * * *", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Staggered != "7 * * * *" {
		t.Errorf("expected '7 * * * *', got %q", res.Staggered)
	}
	if res.OffsetMin != 7 {
		t.Errorf("expected offset 7, got %d", res.OffsetMin)
	}
}

func TestStagger_LiteralMinute(t *testing.T) {
	res, err := Stagger("0 * * * *", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Staggered != "10 * * * *" {
		t.Errorf("expected '10 * * * *', got %q", res.Staggered)
	}
}

func TestStagger_LiteralMinuteWraps(t *testing.T) {
	res, err := Stagger("55 * * * *", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Staggered != "5 * * * *" {
		t.Errorf("expected '5 * * * *', got %q", res.Staggered)
	}
}

func TestStagger_StepExpression(t *testing.T) {
	res, err := Stagger("*/15 * * * *", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// offset 5 % 15 = 5, so minutes: 5, 20, 35, 50
	if res.Staggered != "5,20,35,50 * * * *" {
		t.Errorf("expected '5,20,35,50 * * * *', got %q", res.Staggered)
	}
}

func TestStagger_StepExpressionZeroOffset(t *testing.T) {
	res, err := Stagger("*/20 * * * *", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Staggered != "0,20,40 * * * *" {
		t.Errorf("expected '0,20,40 * * * *', got %q", res.Staggered)
	}
}

func TestStagger_InvalidOffset(t *testing.T) {
	_, err := Stagger("0 * * * *", 60)
	if err == nil {
		t.Fatal("expected error for offset=60")
	}
	_, err = Stagger("0 * * * *", -1)
	if err == nil {
		t.Fatal("expected error for offset=-1")
	}
}

func TestStagger_InvalidExpression(t *testing.T) {
	_, err := Stagger("99 * * * *", 5)
	if err == nil {
		t.Fatal("expected error for invalid cron expression")
	}
}

func TestStagger_WrongFieldCount(t *testing.T) {
	_, err := Stagger("* * *", 5)
	if err == nil {
		t.Fatal("expected error for wrong field count")
	}
}

func TestStagger_ReasonPopulated(t *testing.T) {
	res, err := Stagger("* 2 * * *", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Reason, "wildcard") {
		t.Errorf("expected reason to mention wildcard, got %q", res.Reason)
	}
}
