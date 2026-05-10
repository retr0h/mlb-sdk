// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

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

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// PlayByPlay returns the play-by-play feed for a single game from the v1
// endpoint. For most use cases prefer LiveFeed (v1.1), which sources from
// the same data but is the endpoint MLB Gameday itself uses; this method
// exists for callers who specifically need the older v1 path.
func (c *Client) PlayByPlay(ctx context.Context, gamePk int) ([]Play, error) {
	resp, err := c.raw.GetPlayByPlayWithResponse(ctx, gamePk)
	if err != nil {
		return nil, fmt.Errorf("mlb: playByPlay: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: playByPlay: unexpected status %d", resp.StatusCode())
	}
	return playsFromGenList(resp.JSON200.AllPlays), nil
}

func playsFromGenList(in *[]gen.Play) []Play {
	if in == nil {
		return []Play{}
	}
	out := make([]Play, 0, len(*in))
	for _, p := range *in {
		out = append(out, playFromGen(p))
	}
	return out
}

func playFromGen(p gen.Play) Play {
	out := Play{}
	if p.Result != nil {
		if p.Result.Event != nil {
			out.Event = *p.Result.Event
		}
		if p.Result.EventType != nil {
			out.EventType = EventType(*p.Result.EventType)
		}
		if p.Result.Description != nil {
			out.Description = *p.Result.Description
		}
	}
	if p.About != nil {
		if p.About.Inning != nil {
			out.Inning = *p.About.Inning
		}
		if p.About.HalfInning != nil {
			out.HalfInning = HalfInning(*p.About.HalfInning)
		}
		if p.About.Outs != nil {
			out.Outs = *p.About.Outs
		}
	}
	return out
}
