// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Transactions prints MLB roster / assignment transactions in a date range.
// Run with:
//
//	go run ./examples/transactions [START_DATE] [END_DATE]
//	# default: 2024-07-30 2024-07-31
package main
//
import (
	"context"
	"fmt"
	"os"
	"time"
//
	"github.com/retr0h/mlb-sdk/pkg/mlb"
)
//
const dateFmt = "2006-01-02"
//
func main() {
	startStr, endStr := "2024-07-30", "2024-07-31"
	if len(os.Args) > 2 {
		startStr, endStr = os.Args[1], os.Args[2]
	}
	start, err := time.Parse(dateFmt, startStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "start date:", err)
		os.Exit(2)
	}
	end, err := time.Parse(dateFmt, endStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "end date:", err)
		os.Exit(2)
	}
//
	c := mlb.New()
	tx, err := c.Transactions(context.Background(), mlb.TransactionsQuery{
		StartDate: start,
		EndDate:   end,
		SportID:   1,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "transactions:", err)
		os.Exit(1)
	}
	for _, t := range tx.Transactions {
		fmt.Printf("%s  %-15s %s\n", t.Date, t.TypeDesc, t.Description)
	}
}
