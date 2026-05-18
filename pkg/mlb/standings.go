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
const standingsDateFmt = "2006-01-02"
//
// Standings fetches league standings. q.League is required; 404 is unlikely
// for this endpoint (the API responds with an empty `records` slice when
// nothing matches), but handled below for safety.
//
// Example:
//
//	st, _ := c.Standings(ctx, mlb.StandingsQuery{
//	    League: mlb.NL,
//	    Season: 2026,
//	    StandingsTypes: "regularSeason",
//	})
//	if d := st.Division(204); d != nil {       // NL East
//	    fmt.Println(d.Team(mlb.LAD).Wins)
//	}
func (c *Client) Standings(ctx context.Context, q StandingsQuery) (*Standings, error) {
	if q.League == 0 {
		return nil, fmt.Errorf("mlb: standings: %w: League is required", ErrInvalidQuery)
	}
	params := &gen.GetStandingsParams{LeagueId: q.League.String()}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.StandingsTypes != "" {
		params.StandingsTypes = ptr(q.StandingsTypes)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(standingsDateFmt))
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
//
	resp, err := c.raw.GetStandingsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: standings: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: standings: unexpected status %d", resp.StatusCode())
	}
	return standingsFromGen(resp.JSON200), nil
}
//
func standingsFromGen(r *gen.StandingsResponse) *Standings {
	out := &Standings{}
	if r == nil || r.Records == nil {
		return out
	}
	out.Records = make([]DivisionStandings, 0, len(*r.Records))
	for _, d := range *r.Records {
		out.Records = append(out.Records, divisionStandingsFromGen(d))
	}
	return out
}
//
func divisionStandingsFromGen(d gen.DivisionStandings) DivisionStandings {
	out := DivisionStandings{}
	if d.StandingsType != nil {
		out.StandingsType = *d.StandingsType
	}
	if d.League != nil {
		out.League = refFromGen(d.League)
	}
	if d.Division != nil {
		out.Division = refFromGen(d.Division)
	}
	if d.Sport != nil {
		out.Sport = refFromGen(d.Sport)
	}
	if d.LastUpdated != nil {
		out.LastUpdated = *d.LastUpdated
	}
	if d.TeamRecords != nil {
		out.TeamRecords = make([]TeamRecord, 0, len(*d.TeamRecords))
		for _, tr := range *d.TeamRecords {
			out.TeamRecords = append(out.TeamRecords, teamRecordFromGen(tr))
		}
	}
	return out
}
//
func teamRecordFromGen(t gen.TeamRecord) TeamRecord {
	out := TeamRecord{}
	if t.Team != nil {
		if t.Team.Id != nil {
			out.Team.ID = TeamID(*t.Team.Id)
		}
		if t.Team.Name != nil {
			out.Team.Name = *t.Team.Name
		}
	}
	if t.Streak != nil {
		out.Streak = streakFromGen(*t.Streak)
	}
	if t.Wins != nil {
		out.Wins = *t.Wins
	}
	if t.Losses != nil {
		out.Losses = *t.Losses
	}
	if t.GamesPlayed != nil {
		out.GamesPlayed = *t.GamesPlayed
	}
	if t.RunsScored != nil {
		out.RunsScored = *t.RunsScored
	}
	if t.RunsAllowed != nil {
		out.RunsAllowed = *t.RunsAllowed
	}
	if t.RunDifferential != nil {
		out.RunDifferential = *t.RunDifferential
	}
	if t.WinningPercentage != nil {
		out.WinningPercentage = *t.WinningPercentage
	}
	if t.GamesBack != nil {
		out.GamesBack = *t.GamesBack
	}
	if t.WildCardGamesBack != nil {
		out.WildCardGamesBack = *t.WildCardGamesBack
	}
	if t.DivisionRank != nil {
		out.DivisionRank = *t.DivisionRank
	}
	if t.LeagueRank != nil {
		out.LeagueRank = *t.LeagueRank
	}
	if t.SportRank != nil {
		out.SportRank = *t.SportRank
	}
	if t.EliminationNumber != nil {
		out.EliminationNumber = *t.EliminationNumber
	}
	if t.MagicNumber != nil {
		out.MagicNumber = *t.MagicNumber
	}
	if t.Clinched != nil {
		out.Clinched = *t.Clinched
	}
	if t.DivisionLeader != nil {
		out.DivisionLeader = *t.DivisionLeader
	}
	if t.DivisionChamp != nil {
		out.DivisionChamp = *t.DivisionChamp
	}
	if t.HasWildcard != nil {
		out.HasWildcard = *t.HasWildcard
	}
	if t.Season != nil {
		out.Season = *t.Season
	}
	if t.LastUpdated != nil {
		out.LastUpdated = *t.LastUpdated
	}
	return out
}
//
func refFromGen(r *gen.Ref) Ref {
	out := Ref{}
	if r.Id != nil {
		out.ID = *r.Id
	}
	if r.Link != nil {
		out.Link = *r.Link
	}
	return out
}
//
func streakFromGen(s gen.Streak) Streak {
	out := Streak{}
	if s.StreakCode != nil {
		out.Code = *s.StreakCode
	}
	if s.StreakType != nil {
		out.Type = *s.StreakType
	}
	if s.StreakNumber != nil {
		out.Number = *s.StreakNumber
	}
	return out
}
