package animations

type Animation struct {
	FirstFrame   int
	Last         int
	Step         int     // how many indices do we move per frames
	SpeedInTps   float32 //how many ticks to the next frame
	frameCounter float32
	frame        int
}

func (a *Animation) GetFirstFrame() int {
	return a.FirstFrame
}

func (a *Animation) Update() {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.SpeedInTps
		a.frame += a.Step
		if a.frame > a.Last {
			// loop back to de begining
			a.frame = a.FirstFrame
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}

func NewAnimation(first, last, step int, speed float32) *Animation {
	return &Animation{
		FirstFrame:   first,
		Last:         last,
		Step:         step,
		SpeedInTps:   speed,
		frameCounter: speed,
		frame:        first,
	}
}
