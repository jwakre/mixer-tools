package stopwatch

import (
	"io"
	"log"
	"time"
)

// StopWatch keeps track of a sequence of durations. Use Start and Stop to mark the sections, then
// write the final result with WriteSummary.
type StopWatch struct {
	entries []stopWatchEntry
	t       time.Time
	W       io.Writer
}

type stopWatchEntry struct {
	name string
	d    time.Duration
	used bool
}

func (sw *StopWatch) Start(name string) {
	logger := log.New(sw.W, "", 0)
	if sw.W != nil {
		if len(sw.entries) > 0 {
			logger.Println()

			var d time.Duration
			for _, e := range sw.entries {
				d += e.d
			}
			logger.Printf("=> %s (%s elapsed since start)\n", name,
				d.Truncate(time.Millisecond))
		} else {
			logger.Printf("=> %s\n", name)
		}
	}
	sw.entries = append(sw.entries, stopWatchEntry{name: name})
	sw.t = time.Now()
}

func (sw *StopWatch) Stop() {
	if len(sw.entries) == 0 {
		return
	}
	last := len(sw.entries) - 1
	if sw.entries[last].used {
		sw.entries = append(sw.entries, sw.entries[last])
		last++
	}
	e := &sw.entries[last]
	e.used = true
	e.d = time.Since(sw.t)
}

func (sw *StopWatch) WriteSummary(w io.Writer) {
	logger := log.New(w, "", 0)
	if len(sw.entries) == 0 {
		return
	}
	max := 0
	for _, e := range sw.entries {
		if len(e.name) > max {
			max = len(e.name)
		}
	}
	var sum time.Duration
	logger.Printf("\nTIMINGS\n")
	for _, e := range sw.entries {
		logger.Printf("  %-*s %s\n", max, e.name, e.d.Truncate(time.Millisecond))
		sum += e.d
	}
	logger.Printf("TOTAL: %s\n", sum.Truncate(time.Millisecond))
}
