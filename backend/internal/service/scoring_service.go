package service

import "strings"

type ScoreConfig struct {
	YesScore float64
	NoScore float64
	Weight float64
	IncludeNullInDenominator bool
}

type ScoreResult struct { Numerator, Denominator, Score float64 }

func ComputeBinaryScore(values []string, cfg ScoreConfig) ScoreResult {
	res := ScoreResult{}
	for _, v := range values {
		n := strings.ToLower(strings.TrimSpace(v))
		if n == "" && !cfg.IncludeNullInDenominator { continue }
		res.Denominator += cfg.Weight
		switch n {
		case "yes", "1", "true": res.Numerator += cfg.YesScore * cfg.Weight
		case "no", "0", "false": res.Numerator += cfg.NoScore * cfg.Weight
		default:
			if !cfg.IncludeNullInDenominator { res.Denominator -= cfg.Weight }
		}
	}
	if res.Denominator > 0 { res.Score = res.Numerator / res.Denominator }
	return res
}
