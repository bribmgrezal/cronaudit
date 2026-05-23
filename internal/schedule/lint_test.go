package schedule

import (
	"testing"
)

func TestLint_CleanExpression(t *testing.T) {
	result, err := Lint("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasIssues() {
		t.Errorf("expected no issues, got %v", result.Issues)
	}
}

func TestLint_RedundantStepOne(t *testing.T) {
	result, err := Lint("*/1 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected lint issue for */1")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "minute" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected issue on minute field, got %v", result.Issues)
	}
}

func TestLint_LeadingZero(t *testing.T) {
	result, err := Lint("05 08 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected lint issues for leading zeros")
	}
	count := 0
	for _, issue := range result.Issues {
		if issue.Field == "minute" || issue.Field == "hour" {
			count++
		}
	}
	if count < 2 {
		t.Errorf("expected issues on minute and hour, got %v", result.Issues)
	}
}

func TestLint_EqualRangeEnds(t *testing.T) {
	result, err := Lint("5-5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected lint issue for 5-5 range")
	}
}

func TestLint_BothDOMAndDOW(t *testing.T) {
	result, err := Lint("0 9 15 * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Fatal("expected lint issue for DOM+DOW conflict")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Field == "day-of-month+day-of-week" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected DOM+DOW issue, got %v", result.Issues)
	}
}

func TestLint_InvalidExpression(t *testing.T) {
	_, err := Lint("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestLintResult_HasIssues(t *testing.T) {
	r := LintResult{}
	if r.HasIssues() {
		t.Error("empty result should have no issues")
	}
	r.Issues = append(r.Issues, LintIssue{Field: "minute", Message: "test"})
	if !r.HasIssues() {
		t.Error("result with issues should return true")
	}
}
