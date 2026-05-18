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
const transactionsDateFmt = "2006-01-02"
//
// Transactions fetches roster / assignment transactions. One of the
// following filter forms is required (toddrob99: required_params=[
// ["teamId"], ["playerId"], ["date"], ["startDate", "endDate"]]):
//
//   - q.TeamID set,
//   - q.PlayerID set,
//   - q.On set,
//   - q.StartDate AND q.EndDate set together.
//
// Other combinations return ErrInvalidQuery before any HTTP call.
//
// Example:
//
//	tx, _ := c.Transactions(ctx, mlb.TransactionsQuery{
//	    StartDate: time.Date(2024, 7, 30, 0,0,0,0, time.UTC),
//	    EndDate:   time.Date(2024, 7, 31, 0,0,0,0, time.UTC),
//	})
//	for _, t := range tx.Transactions {
//	    fmt.Println(t.Date, t.TypeDesc, t.Description)
//	}
func (c *Client) Transactions(
	ctx context.Context,
	q TransactionsQuery,
) (*Transactions, error) {
	if err := q.validate(); err != nil {
		return nil, err
	}
//
	params := &gen.GetTransactionsParams{}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.PlayerID != 0 {
		params.PlayerId = ptr(q.PlayerID)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(transactionsDateFmt))
	}
	if !q.StartDate.IsZero() {
		params.StartDate = ptr(q.StartDate.Format(transactionsDateFmt))
	}
	if !q.EndDate.IsZero() {
		params.EndDate = ptr(q.EndDate.Format(transactionsDateFmt))
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetTransactionsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: transactions: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: transactions: unexpected status %d", resp.StatusCode())
	}
	return transactionsFromGen(resp.JSON200), nil
}
//
// validate checks the one-of-required-combos rule. StartDate+EndDate must
// be set together; the other three (TeamID, PlayerID, On) are accepted alone.
func (q TransactionsQuery) validate() error {
	hasTeam := q.TeamID != 0
	hasPlayer := q.PlayerID != 0
	hasOn := !q.On.IsZero()
	hasStart := !q.StartDate.IsZero()
	hasEnd := !q.EndDate.IsZero()
//
	if hasStart != hasEnd {
		return fmt.Errorf(
			"mlb: transactions: %w: StartDate and EndDate must be set together",
			ErrInvalidQuery,
		)
	}
	if !hasTeam && !hasPlayer && !hasOn && !hasStart {
		return fmt.Errorf(
			"mlb: transactions: %w: one of TeamID, PlayerID, On, or StartDate+EndDate is required",
			ErrInvalidQuery,
		)
	}
	return nil
}
//
func transactionsFromGen(r *gen.TransactionsResponse) *Transactions {
	out := &Transactions{}
	if r == nil || r.Transactions == nil {
		return out
	}
	out.Transactions = make([]Transaction, 0, len(*r.Transactions))
	for _, t := range *r.Transactions {
		out.Transactions = append(out.Transactions, transactionFromGen(t))
	}
	return out
}
//
func transactionFromGen(t gen.Transaction) Transaction {
	out := Transaction{}
	if t.Id != nil {
		out.ID = *t.Id
	}
	if t.Person != nil {
		out.Person = personFromGen(*t.Person)
	}
	if t.FromTeam != nil {
		if t.FromTeam.Id != nil {
			out.FromTeam.ID = TeamID(*t.FromTeam.Id)
		}
		if t.FromTeam.Name != nil {
			out.FromTeam.Name = *t.FromTeam.Name
		}
	}
	if t.ToTeam != nil {
		if t.ToTeam.Id != nil {
			out.ToTeam.ID = TeamID(*t.ToTeam.Id)
		}
		if t.ToTeam.Name != nil {
			out.ToTeam.Name = *t.ToTeam.Name
		}
	}
	if t.Date != nil {
		out.Date = *t.Date
	}
	if t.EffectiveDate != nil {
		out.EffectiveDate = *t.EffectiveDate
	}
	if t.ResolutionDate != nil {
		out.ResolutionDate = *t.ResolutionDate
	}
	if t.TypeCode != nil {
		out.TypeCode = *t.TypeCode
	}
	if t.TypeDesc != nil {
		out.TypeDesc = *t.TypeDesc
	}
	if t.Description != nil {
		out.Description = *t.Description
	}
	return out
}
