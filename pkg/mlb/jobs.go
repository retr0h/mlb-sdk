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
// JobsQuery filters a jobs lookup. JobType is required.
type JobsQuery struct {
	JobType string // required: "UMPR", "SCOR", "DCST", …
	SportID int
	Date    string
	Fields  string
}
//
// Jobs fetches staff by job type. q.JobType is required.
func (c *Client) Jobs(ctx context.Context, q JobsQuery) (*Staff, error) {
	if q.JobType == "" {
		return nil, fmt.Errorf("mlb: jobs: %w: JobType is required", ErrInvalidQuery)
	}
	params := &gen.GetJobsParams{JobType: q.JobType}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Date != "" {
		params.Date = ptr(q.Date)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetJobsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: jobs: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: jobs: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
// DatacastersQuery filters a datacasters lookup.
type DatacastersQuery struct {
	SportID int
	Date    string
	Fields  string
}
//
// Datacasters fetches the datacaster roster.
func (c *Client) Datacasters(ctx context.Context, q DatacastersQuery) (*Staff, error) {
	params := &gen.GetJobsDatacastersParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Date != "" {
		params.Date = ptr(q.Date)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetJobsDatacastersWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: datacasters: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: datacasters: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
// OfficialScorersQuery filters an official-scorers lookup.
type OfficialScorersQuery struct {
	Timecode string
	Fields   string
}
//
// OfficialScorers fetches the official scorer roster.
func (c *Client) OfficialScorers(ctx context.Context, q OfficialScorersQuery) (*Staff, error) {
	params := &gen.GetJobsOfficialScorersParams{}
	if q.Timecode != "" {
		params.Timecode = ptr(q.Timecode)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetJobsOfficialScorersWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: officialScorers: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: officialScorers: unexpected status %d", resp.StatusCode())
	}
	return staffFromGen(resp.JSON200), nil
}
//
// PeopleChangesQuery filters a people-changes lookup.
type PeopleChangesQuery struct {
	UpdatedSince string
	Fields       string
}
//
// PeopleChanges fetches recently changed person records.
func (c *Client) PeopleChanges(ctx context.Context, q PeopleChangesQuery) ([]PersonDetail, error) {
	params := &gen.GetPeopleChangesParams{}
	if q.UpdatedSince != "" {
		params.UpdatedSince = ptr(q.UpdatedSince)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetPeopleChangesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: peopleChanges: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: peopleChanges: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
