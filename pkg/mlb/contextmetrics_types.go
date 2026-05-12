// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// ContextMetricsQuery refines a context-metrics lookup. The gamePk path
// parameter is taken as a method argument.
type ContextMetricsQuery struct {
	Timecode string
	Fields   string
}

// ContextMetrics is the typed view of
// /api/v1/game/{gamePk}/contextMetrics.
type ContextMetrics struct {
	Game                         Ref
	LeftFieldSacFlyProbability   map[string]any
	CenterFieldSacFlyProbability map[string]any
	RightFieldSacFlyProbability  map[string]any
	HomeWinProbability           float64
	AwayWinProbability           float64
}
