package schedule

import (
	"testing"
)

func TestSimilarityScore_Identical(t *testing.T) {
	score, err := SimilarityScore("0 * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 1.0 {
		t.Errorf("expected 1.0, got %f", score)
	}
}

func TestSimilarityScore_NoOverlap(t *testing.T) {
	// every hour at :00 vs every hour at :30
	score, err := SimilarityScore("0 * * * *", "30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score != 0.0 {
		t.Errorf("expected 0.0, got %f", score)
	}
}

func TestSimilarityScore_PartialOverlap(t *testing.T) {
	// every 30 min vs every hour
	score, err := SimilarityScore("*/30 * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score <= 0.0 || score >= 1.0 {
		t.Errorf("expected partial overlap score, got %f", score)
	}
}

func TestSimilarityScore_InvalidExprA(t *testing.T) {
	_, err := SimilarityScore("bad expr", "0 * * * *")
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestSimilarityScore_InvalidExprB(t *testing.T) {
	_, err := SimilarityScore("0 * * * *", "bad expr")
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestGroupSimilar_AllSame(t *testing.T) {
	exprs := []string{"0 * * * *", "0 * * * *", "0 * * * *"}
	groups := GroupSimilar(exprs, 1.0)
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0]) != 3 {
		t.Errorf("expected group of 3, got %d", len(groups[0]))
	}
}

func TestGroupSimilar_AllDifferent(t *testing.T) {
	exprs := []string{"0 * * * *", "30 * * * *", "15 * * * *"}
	groups := GroupSimilar(exprs, 1.0)
	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}
}

func TestGroupSimilar_EmptyInput(t *testing.T) {
	groups := GroupSimilar([]string{}, 0.5)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}
