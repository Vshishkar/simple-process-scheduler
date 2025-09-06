package main

import "time"

func averageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range durations {
		total += t
	}
	return total / time.Duration(len(durations))
}
