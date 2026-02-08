package ui

// FirstFrameTracker handles the "skip first frame" pattern used by UI states
// that need to ignore the input event which opened them.
type FirstFrameTracker struct {
	firstFrame bool
}

func NewFirstFrameTracker() FirstFrameTracker {
	return FirstFrameTracker{firstFrame: true}
}

func (f *FirstFrameTracker) IsFirstFrame() bool {
	return f.firstFrame
}

func (f *FirstFrameTracker) ClearFirstFrame() {
	f.firstFrame = false
}

func (f *FirstFrameTracker) ResetFirstFrame() {
	f.firstFrame = true
}
