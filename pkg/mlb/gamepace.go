// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import (
	"context"
	"fmt"
//
	"github.com/retr0h/mlb-sdk/internal/gen"
)
//
// GamePace fetches pace-of-play statistics. q.Season is required (toddrob99:
// required_params=[["season"]]).
//
// Example:
//
//	gp, _ := c.GamePace(ctx, mlb.GamePaceQuery{Season: 2024, SportID: 1})
//	fmt.Println("time per game:", gp.TimePerGame, "pitches:", gp.TotalPitches)
func (c *Client) GamePace(ctx context.Context, q GamePaceQuery) (*GamePace, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: gamePace: %w: Season is required", ErrInvalidQuery)
	}
//
	params := &gen.GetGamePaceParams{Season: q.Season}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.TeamIDs != "" {
		params.TeamIds = ptr(q.TeamIDs)
	}
	if q.LeagueIDs != "" {
		params.LeagueIds = ptr(q.LeagueIDs)
	}
	if q.LeagueListID != "" {
		params.LeagueListId = ptr(q.LeagueListID)
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.StartDate != "" {
		params.StartDate = ptr(q.StartDate)
	}
	if q.EndDate != "" {
		params.EndDate = ptr(q.EndDate)
	}
	if q.VenueIDs != "" {
		params.VenueIds = ptr(q.VenueIDs)
	}
	if q.OrgType != "" {
		params.OrgType = ptr(q.OrgType)
	}
	if q.IncludeChildren {
		params.IncludeChildren = ptr(q.IncludeChildren)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetGamePaceWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gamePace: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gamePace: unexpected status %d", resp.StatusCode())
	}
	gp := gamePaceFromResponse(resp.JSON200)
	if gp == nil {
		return nil, ErrNotFound
	}
	return gp, nil
}
//
func gamePaceFromResponse(r *gen.GamePaceResponse) *GamePace {
	if r == nil || r.Sports == nil || len(*r.Sports) == 0 {
		return nil
	}
	gp := gamePaceFromGen((*r.Sports)[0])
	return &gp
}
//
func gamePaceFromGen(g gen.GamePaceData) GamePace {
	out := GamePace{}
	if g.Season != nil {
		out.Season = *g.Season
	}
	if g.Sport != nil {
		out.Sport = refFromGen(g.Sport)
	}
	if g.HitsPer9Inn != nil {
		out.HitsPer9Inn = *g.HitsPer9Inn
	}
	if g.RunsPer9Inn != nil {
		out.RunsPer9Inn = *g.RunsPer9Inn
	}
	if g.PitchesPer9Inn != nil {
		out.PitchesPer9Inn = *g.PitchesPer9Inn
	}
	if g.PlateAppearancesPer9Inn != nil {
		out.PlateAppearancesPer9Inn = *g.PlateAppearancesPer9Inn
	}
	if g.HitsPerGame != nil {
		out.HitsPerGame = *g.HitsPerGame
	}
	if g.RunsPerGame != nil {
		out.RunsPerGame = *g.RunsPerGame
	}
	if g.InningsPlayedPerGame != nil {
		out.InningsPlayedPerGame = *g.InningsPlayedPerGame
	}
	if g.PitchesPerGame != nil {
		out.PitchesPerGame = *g.PitchesPerGame
	}
	if g.PitchersPerGame != nil {
		out.PitchersPerGame = *g.PitchersPerGame
	}
	if g.PlateAppearancesPerGame != nil {
		out.PlateAppearancesPerGame = *g.PlateAppearancesPerGame
	}
	if g.HitsPerRun != nil {
		out.HitsPerRun = *g.HitsPerRun
	}
	if g.PitchesPerPitcher != nil {
		out.PitchesPerPitcher = *g.PitchesPerPitcher
	}
	if g.TotalGameTime != nil {
		out.TotalGameTime = *g.TotalGameTime
	}
	if g.TotalInningsPlayed != nil {
		out.TotalInningsPlayed = *g.TotalInningsPlayed
	}
	if g.TotalHits != nil {
		out.TotalHits = *g.TotalHits
	}
	if g.TotalRuns != nil {
		out.TotalRuns = *g.TotalRuns
	}
	if g.TotalPlateAppearances != nil {
		out.TotalPlateAppearances = *g.TotalPlateAppearances
	}
	if g.TotalPitchers != nil {
		out.TotalPitchers = *g.TotalPitchers
	}
	if g.TotalPitches != nil {
		out.TotalPitches = *g.TotalPitches
	}
	if g.TotalGames != nil {
		out.TotalGames = *g.TotalGames
	}
	if g.Total7InnGames != nil {
		out.Total7InnGames = *g.Total7InnGames
	}
	if g.Total9InnGames != nil {
		out.Total9InnGames = *g.Total9InnGames
	}
	if g.Total9InnGamesCompletedEarly != nil {
		out.Total9InnGamesCompletedEarly = *g.Total9InnGamesCompletedEarly
	}
	if g.Total9InnGamesScheduled != nil {
		out.Total9InnGamesScheduled = *g.Total9InnGamesScheduled
	}
	if g.Total9InnGamesWithoutExtraInn != nil {
		out.Total9InnGamesWithoutExtraInn = *g.Total9InnGamesWithoutExtraInn
	}
	if g.TotalExtraInnGames != nil {
		out.TotalExtraInnGames = *g.TotalExtraInnGames
	}
	if g.TotalExtraInnTime != nil {
		out.TotalExtraInnTime = *g.TotalExtraInnTime
	}
	if g.TimePerGame != nil {
		out.TimePerGame = *g.TimePerGame
	}
	if g.TimePerPitch != nil {
		out.TimePerPitch = *g.TimePerPitch
	}
	if g.TimePerHit != nil {
		out.TimePerHit = *g.TimePerHit
	}
	if g.TimePerRun != nil {
		out.TimePerRun = *g.TimePerRun
	}
	if g.TimePerPlateAppearance != nil {
		out.TimePerPlateAppearance = *g.TimePerPlateAppearance
	}
	if g.TimePer9Inn != nil {
		out.TimePer9Inn = *g.TimePer9Inn
	}
	if g.TimePer77PlateAppearances != nil {
		out.TimePer77PlateAppearances = *g.TimePer77PlateAppearances
	}
	if g.TimePer7InnGameWithoutExtraInn != nil {
		out.TimePer7InnGameWithoutExtraInn = *g.TimePer7InnGameWithoutExtraInn
	}
	return out
}
