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
// GameChanges fetches games modified since a timestamp. q.UpdatedSince is
// required (toddrob99: required_params=[["updatedSince"]]).
//
// Example:
//
//	gc, _ := c.GameChanges(ctx, mlb.GameChangesQuery{
//	    UpdatedSince: "2024-09-07T00:00:00", SportID: 1,
//	})
//	fmt.Println("changed games:", gc.TotalGames)
func (c *Client) GameChanges(
	ctx context.Context,
	q GameChangesQuery,
) (*GameChanges, error) {
	if q.UpdatedSince == "" {
		return nil, fmt.Errorf(
			"mlb: gameChanges: %w: UpdatedSince is required",
			ErrInvalidQuery,
		)
	}
//
	params := &gen.GetGameChangesParams{UpdatedSince: q.UpdatedSince}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetGameChangesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameChanges: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameChanges: unexpected status %d", resp.StatusCode())
	}
	return gameChangesFromGen(resp.JSON200), nil
}
//
func gameChangesFromGen(r *gen.GameChangesResponse) *GameChanges {
	out := &GameChanges{}
	if r == nil {
		return out
	}
	if r.TotalItems != nil {
		out.TotalItems = *r.TotalItems
	}
	if r.TotalEvents != nil {
		out.TotalEvents = *r.TotalEvents
	}
	if r.TotalGames != nil {
		out.TotalGames = *r.TotalGames
	}
	if r.TotalGamesInProgress != nil {
		out.TotalGamesInProgress = *r.TotalGamesInProgress
	}
	if r.Dates != nil {
		out.Dates = make([]GameChangesDate, 0, len(*r.Dates))
		for _, d := range *r.Dates {
			out.Dates = append(out.Dates, gameChangesDateFromGen(d))
		}
	}
	return out
}
//
func gameChangesDateFromGen(d gen.ScheduleDate) GameChangesDate {
	out := GameChangesDate{}
	if d.Date != nil {
		out.Date = *d.Date
	}
	if d.Games != nil {
		out.Games = make([]Game, 0, len(*d.Games))
		for _, g := range *d.Games {
			out.Games = append(out.Games, gameFromGen(g))
		}
	}
	return out
}
