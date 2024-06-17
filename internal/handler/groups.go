package handler

import (
	"context"
)

// ReceiveGroupsFromMatchmaker - receive groups with players from matchmaker and prints it to stdout with calculated metrics.
func (mh *MatchmakingHandler) ReceiveGroupsFromMatchmaker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case group, ok := <-mh.matchmaker.ReceiveGroups():
				if ok {
					group.CalculateMetrics()
					group.Print()
				}
			}
		}
	}()
}
