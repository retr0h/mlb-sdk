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
// Venue fetches a single venue by id. The MLB API's `/api/v1/venues` endpoint
// also supports a comma-separated `venueIds` query parameter, but this SDK
// wraps the path-based single-object lookup form. Hydrate "location",
// "fieldInfo", and "timezone" via q.Hydrate to populate the optional sub-
// blocks.
//
// Example:
//
//	v, _ := c.Venue(ctx, 22, mlb.VenueQuery{Hydrate: "location,fieldInfo,timezone"})
//	fmt.Println(v.Name, v.Location.City, v.FieldInfo.Capacity)
func (c *Client) Venue(ctx context.Context, venueID int, q VenueQuery) (*Venue, error) {
	params := &gen.GetVenueParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetVenueWithResponse(ctx, venueID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: venue: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: venue: unexpected status %d", resp.StatusCode())
	}
	v := venueFromResponse(resp.JSON200)
	if v == nil {
		return nil, ErrNotFound
	}
	return v, nil
}
//
// venueFromResponse picks the first (and typically only) venue from the
// generated response. Returns nil when the venues slice is empty — the
// caller maps that to ErrNotFound.
func venueFromResponse(r *gen.VenueResponse) *Venue {
	if r == nil || r.Venues == nil || len(*r.Venues) == 0 {
		return nil
	}
	v := venueFromGen((*r.Venues)[0])
	return &v
}
//
func venueFromGen(v gen.Venue) Venue {
	out := Venue{}
	if v.Id != nil {
		out.ID = *v.Id
	}
	if v.Name != nil {
		out.Name = *v.Name
	}
	if v.Link != nil {
		out.Link = *v.Link
	}
	if v.Active != nil {
		out.Active = *v.Active
	}
	if v.Season != nil {
		out.Season = *v.Season
	}
	if v.Location != nil {
		out.Location = venueLocationFromGen(*v.Location)
	}
	if v.TimeZone != nil {
		out.TimeZone = venueTimeZoneFromGen(*v.TimeZone)
	}
	if v.FieldInfo != nil {
		out.FieldInfo = venueFieldInfoFromGen(*v.FieldInfo)
	}
	return out
}
//
func venueLocationFromGen(l gen.VenueLocation) VenueLocation {
	out := VenueLocation{}
	if l.Address1 != nil {
		out.Address1 = *l.Address1
	}
	if l.Address2 != nil {
		out.Address2 = *l.Address2
	}
	if l.City != nil {
		out.City = *l.City
	}
	if l.State != nil {
		out.State = *l.State
	}
	if l.StateAbbrev != nil {
		out.StateAbbrev = *l.StateAbbrev
	}
	if l.PostalCode != nil {
		out.PostalCode = *l.PostalCode
	}
	if l.DefaultCoordinates != nil {
		out.DefaultCoordinates = venueCoordinatesFromGen(*l.DefaultCoordinates)
	}
	if l.AzimuthAngle != nil {
		out.AzimuthAngle = *l.AzimuthAngle
	}
	if l.Elevation != nil {
		out.Elevation = *l.Elevation
	}
	if l.Country != nil {
		out.Country = *l.Country
	}
	if l.Phone != nil {
		out.Phone = *l.Phone
	}
	return out
}
//
func venueCoordinatesFromGen(c gen.VenueCoordinates) VenueCoordinates {
	out := VenueCoordinates{}
	if c.Latitude != nil {
		out.Latitude = *c.Latitude
	}
	if c.Longitude != nil {
		out.Longitude = *c.Longitude
	}
	return out
}
//
func venueTimeZoneFromGen(t gen.VenueTimeZone) VenueTimeZone {
	out := VenueTimeZone{}
	if t.Tz != nil {
		out.TZ = *t.Tz
	}
	if t.Id != nil {
		out.ID = *t.Id
	}
	if t.Offset != nil {
		out.Offset = *t.Offset
	}
	if t.OffsetAtGameTime != nil {
		out.OffsetAtGameTime = *t.OffsetAtGameTime
	}
	return out
}
//
func venueFieldInfoFromGen(f gen.VenueFieldInfo) VenueFieldInfo {
	out := VenueFieldInfo{}
	if f.Capacity != nil {
		out.Capacity = *f.Capacity
	}
	if f.TurfType != nil {
		out.TurfType = *f.TurfType
	}
	if f.RoofType != nil {
		out.RoofType = *f.RoofType
	}
	if f.LeftLine != nil {
		out.LeftLine = *f.LeftLine
	}
	if f.LeftCenter != nil {
		out.LeftCenter = *f.LeftCenter
	}
	if f.Center != nil {
		out.Center = *f.Center
	}
	if f.RightCenter != nil {
		out.RightCenter = *f.RightCenter
	}
	if f.RightLine != nil {
		out.RightLine = *f.RightLine
	}
	return out
}
