// Copyright (c) 2026 John Dewey
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package mlb

import (
	"context"
	"fmt"
)

// LiveFeed returns the live game feed for a single game from the v1.1
// /feed/live endpoint. Today this surface returns just the play list
// (`liveData.plays.allPlays`), which is the same shape PlayByPlay returns
// — but sourced from the v1.1 endpoint MLB Gameday itself uses, so this
// is the recommended entry point.
//
// Future phases will broaden the return type to also surface gameData,
// boxscore, and decisions; this method's signature is therefore expected
// to evolve.
func (c *Client) LiveFeed(ctx context.Context, gamePk int) ([]Play, error) {
	resp, err := c.raw.GetLiveFeedWithResponse(ctx, gamePk)
	if err != nil {
		return nil, fmt.Errorf("mlb: liveFeed: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: liveFeed: unexpected status %d", resp.StatusCode())
	}
	if resp.JSON200.LiveData == nil || resp.JSON200.LiveData.Plays == nil {
		return []Play{}, nil
	}
	return playsFromGenList(resp.JSON200.LiveData.Plays.AllPlays), nil
}
