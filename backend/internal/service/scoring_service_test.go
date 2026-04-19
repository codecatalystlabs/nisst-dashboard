package service

import "testing"

func TestComputeBinaryScore(t *testing.T) {
	cfg := ScoreConfig{YesScore: 1, NoScore: 0, Weight: 1}
	res := ComputeBinaryScore([]string{"Yes", "No", "Yes", ""}, cfg)
	if res.Denominator != 3 { t.Fatalf("expected denominator=3, got %v", res.Denominator) }
	if res.Score < 0.66 || res.Score > 0.67 { t.Fatalf("expected ~0.6667, got %v", res.Score) }
}
