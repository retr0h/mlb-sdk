// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// ContextMetrics prints win probability for a game. Run with:
//
//	go run ./examples/contextmetrics [GAME_PK]   # default: 745455
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	gamePk := 745455
	if len(os.Args) > 1 {
		n, _ := strconv.Atoi(os.Args[1])
		gamePk = n
	}

	c := mlb.New()
	cm, err := c.ContextMetrics(context.Background(), gamePk, mlb.ContextMetricsQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "contextMetrics:", err)
		os.Exit(1)
	}
	fmt.Printf("game %d — home %.1f%% away %.1f%%\n",
		cm.Game.ID, cm.HomeWinProbability, cm.AwayWinProbability)
}
