// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// VenueQuery refines a venue lookup. The venueId path parameter is taken as
// a method argument; everything in this struct is optional.
type VenueQuery struct {
	// Season constrains the venue to a specific season's metadata (e.g. when
	// a stadium changes name or capacity year-over-year).
	Season int
//
	// Hydrate is a comma-separated free-form hydrate string — common values
	// are "location", "fieldInfo", "timezone". The MLB API's hydrate
	// vocabulary grows; we don't constrain it.
	Hydrate string
//
	// Fields restricts the response to a comma-separated field projection.
	Fields string
}
//
// Venue is the typed view of a single venue from /api/v1/venues/{venueId}.
// Location, TimeZone and FieldInfo are only populated when hydrated.
type Venue struct {
	ID        int
	Name      string
	Link      string
	Active    bool
	Season    string
	Location  VenueLocation
	TimeZone  VenueTimeZone
	FieldInfo VenueFieldInfo
}
//
// VenueLocation is the venue's postal address and geographic placement. The
// MLB API only includes it when the request hydrates `location`.
type VenueLocation struct {
	Address1           string
	Address2           string
	City               string
	State              string
	StateAbbrev        string
	PostalCode         string
	DefaultCoordinates VenueCoordinates
	AzimuthAngle       float64
	Elevation          int
	Country            string
	Phone              string
}
//
// VenueCoordinates is the lat/long pair of the venue's default coordinates.
type VenueCoordinates struct {
	Latitude  float64
	Longitude float64
}
//
// VenueTimeZone is the venue's local time-zone metadata. Only populated when
// the request hydrates `timezone`.
type VenueTimeZone struct {
	TZ               string // e.g. "PDT"
	ID               string // IANA id, e.g. "America/Los_Angeles"
	Offset           int    // hours from UTC
	OffsetAtGameTime int
}
//
// VenueFieldInfo is field dimensions and surface details. Only populated when
// the request hydrates `fieldInfo`.
type VenueFieldInfo struct {
	Capacity    int
	TurfType    string
	RoofType    string
	LeftLine    int
	LeftCenter  int
	Center      int
	RightCenter int
	RightLine   int
}
