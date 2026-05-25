package schedule

import (
	"testing"
)

func TestScoreOverlap_Identical(t *testing.T) {
	score, err := ScoreOverlap("0 * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.ScorePercent != 100.0 {
		t.Errorf("expected 100%% overlap, got %.2f", score.ScorePercent)
	}
	if score.SharedCount == 0 {
		t.Error("expected non-zero shared count for identical schedules")
	}
}

func TestScoreOverlap_NoOverlap(t *testing.T) {
	// minute 0 every hour vs minute 30 every hour — never overlap
	score, err := ScoreOverlap("0 * * * *", "30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.ScorePercent != 0.0 {
		t.Errorf("expected 0%% overlap, got %.2f", score.ScorePercent)
	}
	if score.SharedCount != 0 {
		t.Errorf("expected 0 shared minutes, got %d", score.SharedCount)
	}
}

func TestScoreOverlap_PartialOverlap(t *testing.T) {
	// every minute vs every hour at minute 0 — 24 shared out of 1440+24-24
	score, err := ScoreOverlap("* * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.SharedCount != 24 {
		t.Errorf("expected 24 shared minutes, got %d", score.SharedCount)
	}
	if score.ScorePercent <= 0 || score.ScorePercent >= 100 {
		t.Errorf("expected partial overlap percent, got %.2f", score.ScorePercent)
	}
}

func TestScoreOverlap_InvalidExprA(t *testing.T) {
	_, err := ScoreOverlap("bad expr", "0 * * * *")
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestScoreOverlap_InvalidExprB(t *testing.T) {
	_, err := ScoreOverlap("0 * * * *", "bad expr")
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestRankOverlaps_SortedDescending(t *testing.T) {
	exprs := []string{"0 * * * *", "0 * * * *", "30 * * * *"}
	results, err := RankOverlaps(exprs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].ScorePercent > results[i-1].ScorePercent {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
	// first result should be the identical pair
	if results[0].ScorePercent != 100.0 {
		t.Errorf("expected first pair to be 100%% overlap, got %.2f", results[0].ScorePercent)
	}
}

func TestRankOverlaps_InvalidExpr(t *testing.T) {
	_, err := RankOverlaps([]string{"0 * * * *", "not valid"})
	if err == nil {
		t.Error("expected error for invalid expression in list")
	}
}
