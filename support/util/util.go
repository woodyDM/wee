package util

import "time"

type TimeWatcher struct {
	start time.Time
	end   time.Time
}

func StartTimeWatch() *TimeWatcher {
	watcher := new(TimeWatcher)
	watcher.Start()
	return watcher
}

func (w *TimeWatcher) String() string {
	duration := w.end.Sub(w.start)
	return duration.String()
}

func (w *TimeWatcher) Start() {
	w.start = time.Now()
}

func (w *TimeWatcher) End() {
	w.end = time.Now()
}
