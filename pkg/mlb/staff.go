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
const staffDateFmt = "2006-01-02"
//
// Coaches fetches the coaching staff for a team.
func (c *Client) Coaches(ctx context.Context, teamID int, q CoachesQuery) (*Staff, error) {
	params := &gen.GetTeamCoachesParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(staffDateFmt))
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetTeamCoachesWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: coaches: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: coaches: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
// Personnel fetches front-office personnel for a team.
func (c *Client) Personnel(ctx context.Context, teamID int, q PersonnelQuery) (*Staff, error) {
	params := &gen.GetTeamPersonnelParams{}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(staffDateFmt))
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetTeamPersonnelWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: personnel: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: personnel: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
// Umpires fetches the umpire roster.
func (c *Client) Umpires(ctx context.Context, q UmpiresQuery) (*Staff, error) {
	params := &gen.GetJobsUmpiresParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(staffDateFmt))
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetJobsUmpiresWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: umpires: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: umpires: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
func staffFromGen(r *gen.StaffResponse) *Staff {
	out := &Staff{}
	if r == nil {
		return out
	}
	if r.Link != nil {
		out.Link = *r.Link
	}
	if r.TeamId != nil {
		out.TeamID = *r.TeamId
	}
	if r.RosterType != nil {
		out.RosterType = *r.RosterType
	}
	if r.Roster != nil {
		out.Roster = make([]StaffEntry, 0, len(*r.Roster))
		for _, e := range *r.Roster {
			out.Roster = append(out.Roster, staffEntryFromGen(e))
		}
	}
	return out
}
//
func staffEntryFromGen(e gen.StaffEntry) StaffEntry {
	out := StaffEntry{}
	if e.Person != nil {
		out.Person = personFromGen(*e.Person)
	}
	if e.JerseyNumber != nil {
		out.JerseyNumber = *e.JerseyNumber
	}
	if e.Job != nil {
		out.Job = *e.Job
	}
	if e.JobId != nil {
		out.JobID = *e.JobId
	}
	if e.Title != nil {
		out.Title = *e.Title
	}
	return out
}
