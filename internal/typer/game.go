// Package typer implements the typing game logic.
package typer

import "time"

// Result holds the final statistics for a completed game.
type Result struct {
	NetWPM   int
	RawWPM   int
	Accuracy float64
	Elapsed  time.Duration
}

// Score computes game statistics from the target text, user input, and elapsed time.
func Score(target, input string, elapsed time.Duration) Result {
	t := []rune(target)
	in := []rune(input)

	correct := 0
	for i, r := range in {
		if i < len(t) && r == t[i] {
			correct++
		}
	}

	minutes := elapsed.Minutes()
	if minutes < 0.001 {
		minutes = 0.001
	}

	accuracy := 0.0
	if len(in) > 0 {
		accuracy = float64(correct) / float64(len(in)) * 100
	}

	return Result{
		NetWPM:   int(float64(correct) / 5 / minutes),
		RawWPM:   int(float64(len(in)) / 5 / minutes),
		Accuracy: accuracy,
		Elapsed:  elapsed,
	}
}
