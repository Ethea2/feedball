package animation

type ThrowPhase uint8

const (
	WindupPhase ThrowPhase = iota
	ReleasePhase
)

type ThrowAnimation struct {
	windup  *Animation
	release *Animation
	phase   ThrowPhase
	done    bool
}

func NewThrowAnimation(
	windupFirst, windupLast int,
	releaseFirst, releaseLast int,
	speed float32,
) *ThrowAnimation {
	return &ThrowAnimation{
		windup:  NewAnimation(windupFirst, windupLast, 1, speed),
		release: NewAnimation(releaseFirst, releaseLast, 1, speed),
		phase:   WindupPhase,
		done:    false,
	}
}

func (t *ThrowAnimation) IsOnLastWindupFrame() bool {
	return t.phase == WindupPhase && t.windup.frame == t.windup.Last
}

func (t *ThrowAnimation) IsDone() bool {
	return t.done
}

func (t *ThrowAnimation) SetRelease(first, last int) {
	t.release.frame = first
	t.release.First = first
	t.release.Last = last
	t.release.frameCounter = t.release.SpeedInTps
}

func (t *ThrowAnimation) Frame() int {
	if t.phase == WindupPhase {
		return t.windup.frame
	}
	return t.release.frame
}

func (t *ThrowAnimation) Update() {
	if t.done {
		return
	}

	if t.phase == WindupPhase {
		t.windup.frameCounter -= 1.0
		if t.windup.frameCounter < 0.0 {
			t.windup.frameCounter = t.windup.SpeedInTps
			if t.windup.frame == t.windup.Last {
				t.phase = ReleasePhase
			} else {
				t.windup.frame += t.windup.Step
			}
		}
	} else {
		t.release.frameCounter -= 1.0
		if t.release.frameCounter < 0.0 {
			t.release.frameCounter = t.release.SpeedInTps
			if t.release.frame == t.release.Last {
				t.done = true
			} else {
				t.release.frame += t.release.Step
			}
		}
	}
}
