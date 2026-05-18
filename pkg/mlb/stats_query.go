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
// StatsQuery filters a league-wide stats lookup. Stats and Group are required.
type StatsQuery struct {
	Stats      string // required: "season", "career", …
	Group      string // required: "hitting", "pitching", "fielding"
	Season     int
	SportIDs   string
	GameType   string
	PlayerPool string
	Position   string
	TeamID     int
	LeagueID   int
	PersonID   int
	Limit      int
	Offset     int
	SortStat   string
	Order      string
	Metrics    string
	StartDate  string
	EndDate    string
	Hydrate    string
	Fields     string
}
//
// Stats fetches league-wide individual player stats. q.Stats and q.Group
// are required.
func (c *Client) Stats(ctx context.Context, q StatsQuery) (*TeamStats, error) {
	if q.Stats == "" || q.Group == "" {
		return nil, fmt.Errorf(
			"mlb: stats: %w: Stats and Group are both required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetStatsParams{Stats: q.Stats, Group: q.Group}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.SportIDs != "" {
		params.SportIds = ptr(q.SportIDs)
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.PlayerPool != "" {
		params.PlayerPool = ptr(q.PlayerPool)
	}
	if q.Position != "" {
		params.Position = ptr(q.Position)
	}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.PersonID != 0 {
		params.PersonId = ptr(q.PersonID)
	}
	if q.Limit != 0 {
		params.Limit = ptr(q.Limit)
	}
	if q.Offset != 0 {
		params.Offset = ptr(q.Offset)
	}
	if q.SortStat != "" {
		params.SortStat = ptr(q.SortStat)
	}
	if q.Order != "" {
		params.Order = ptr(q.Order)
	}
	if q.Metrics != "" {
		params.Metrics = ptr(q.Metrics)
	}
	if q.StartDate != "" {
		params.StartDate = ptr(q.StartDate)
	}
	if q.EndDate != "" {
		params.EndDate = ptr(q.EndDate)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetStatsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: stats: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: stats: unexpected status %d", resp.StatusCode())
	}
	return teamStatsFromGen(resp.JSON200), nil
}
//
// SchedulePostseasonSeries fetches postseason series data.
func (c *Client) SchedulePostseasonSeries(
	ctx context.Context,
	q SchedulePostseasonQuery,
) ([]Game, error) {
	params := &gen.GetSchedulePostseasonSeriesParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.GameTypes != "" {
		params.GameTypes = ptr(q.GameTypes)
	}
	if q.SeriesNumber != 0 {
		params.SeriesNumber = ptr(q.SeriesNumber)
	}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetSchedulePostseasonSeriesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: schedulePostseasonSeries: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf(
			"mlb: schedulePostseasonSeries: unexpected status %d",
			resp.StatusCode(),
		)
	}
	var games []Game
	if resp.JSON200.Dates != nil {
		for _, d := range *resp.JSON200.Dates {
			if d.Games != nil {
				for _, g := range *d.Games {
					games = append(games, gameFromGen(g))
				}
			}
		}
	}
	return games, nil
}
//
// AllStarBallotQuery filters an All-Star ballot lookup. Season is required.
type AllStarBallotQuery struct {
	Season int // required
	Fields string
}
//
// AllStarBallot fetches the All-Star ballot for a league.
func (c *Client) AllStarBallot(
	ctx context.Context,
	leagueID int,
	q AllStarBallotQuery,
) ([]PersonDetail, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: allStarBallot: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetAllStarBallotParams{Season: q.Season}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
	resp, err := c.raw.GetAllStarBallotWithResponse(ctx, leagueID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: allStarBallot: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: allStarBallot: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
//
// AllStarFinalVote fetches the All-Star final vote for a league.
func (c *Client) AllStarFinalVote(
	ctx context.Context,
	leagueID int,
	q AllStarBallotQuery,
) ([]PersonDetail, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: allStarFinalVote: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetAllStarFinalVoteParams{Season: q.Season}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
	resp, err := c.raw.GetAllStarFinalVoteWithResponse(ctx, leagueID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: allStarFinalVote: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: allStarFinalVote: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
//
// AllStarWriteIns fetches the All-Star write-in votes for a league.
func (c *Client) AllStarWriteIns(
	ctx context.Context,
	leagueID int,
	q AllStarBallotQuery,
) ([]PersonDetail, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: allStarWriteIns: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetAllStarWriteInsParams{Season: q.Season}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
	resp, err := c.raw.GetAllStarWriteInsWithResponse(ctx, leagueID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: allStarWriteIns: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: allStarWriteIns: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
