package models

import (
	"log"
	"math"
	"math/rand"
	"time"
)

const MicrosecondsToX = 40_000
const StartingTubeSpacing = 150
const SpacingStepDown = 50
const MINIMUM_TUBE_SPACING = 12
const REDUCTION_SCALE = 3
const PIPE_DIST_FROM_EDGES = 6
const PIPE_HOLE_SIZE = 15

type Pipes struct {
	context *Context

	totalPipes int

	elapsedTime     int64
	currentStep     int64
	lastPipeCreated int64

	Pipes []*Pipe
}

type Pipe struct {
	x      int
	offset int

	opening int
	top     int
	bottom  int

	context *Context
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func newPipe(context *Context, startingX int) *Pipe {
    _, height := context.Terminal.GetFixedBounds()
	p := &Pipe{
		context: context,
		x:       startingX,
	}

	p.sizeUpPipe(height)

	return p
}

func (p *Pipe) CreateRender(size int) (*Point, [][]byte) {
	_, height := p.context.Terminal.GetFixedBounds()

    height *= 0x1 << size

	point := &Point{
		X: float64(p.x),
		Y: float64(0),
	}

	log.Printf("Pipe#CreateRender(%v): %v %v %v", size, p.top, p.bottom, height)

	display := make([][]byte, height)
	for i := 0; i < height; i += 1 {
		if i <= p.top * (size + 1) || i >= p.bottom * (size + 1) {
			display[i] = getPipeDesign(size)
		} else {
			display[i] = getEmptySpace(size)
		}
	}

	return point, display
}

func getPipeDesign(scale int) []byte {
	switch scale {
	case 1:
		return []byte{'|', ' ', ' ', '|'}
	case 2:
		return []byte{'|', '.', ' ', ' ', ' ', ' ', '.', '|'}
	}

	return []byte{'x'}
}

func getEmptySpace(scale int) []byte {
	out := []byte{}
	for i := 0; i < int(0x1<<scale); i++ {
		out = append(out, ' ')
	}

	return out
}

func (p *Pipe) sizeUpPipe(height int) {

	opening := randInt(PIPE_DIST_FROM_EDGES, height-PIPE_DIST_FROM_EDGES)
	top := opening - PIPE_HOLE_SIZE/2
	bottom := opening + PIPE_DIST_FROM_EDGES - PIPE_DIST_FROM_EDGES/2

	if top < PIPE_DIST_FROM_EDGES {
		diff := PIPE_DIST_FROM_EDGES - top
		bottom += diff
		top = PIPE_DIST_FROM_EDGES
	} else if bottom > height-PIPE_DIST_FROM_EDGES {
		diff := bottom - (height - PIPE_DIST_FROM_EDGES)
		top -= diff
		bottom = height - PIPE_DIST_FROM_EDGES
	}

    p.opening = opening
    p.top = top
    p.bottom = bottom;
}

func NewPipes(context *Context) *Pipes {
	return &Pipes{
		lastPipeCreated: 0,
		elapsedTime:     0,
		currentStep:     0,
		Pipes:           []*Pipe{},
		totalPipes:      0,
		context:         context,
	}
}

func min(one float64, two int) float64 {
	two_f := float64(two)
	if two_f > one {
		return one
	}
	return two_f
}

func maxInt(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (p *Pipes) canCreatePipe() bool {
	if p.lastPipeCreated == 0 {
		return true
	}

	pipeCount := 1
	takenSteps := p.elapsedTime / MicrosecondsToX

	for {
		scaledReduce := pipeCount / REDUCTION_SCALE

		currentStepsRequired := maxInt(
			int64(150*math.Pow(float64(scaledReduce+1), -.71)),
			MINIMUM_TUBE_SPACING,
		)

		if takenSteps < currentStepsRequired {
			break
		}

		pipeCount += 1
		takenSteps -= currentStepsRequired
	}

	return p.totalPipes < pipeCount
}

func (p *Pipes) Update(delta time.Duration) {
	width, _ := p.context.Terminal.GetFixedBounds()
	p.elapsedTime += delta.Microseconds()
	if p.canCreatePipe() {
		pipe := newPipe(p.context, width-3)
		p.Pipes = append(p.Pipes, pipe)
		p.lastPipeCreated = p.elapsedTime
		p.totalPipes += 1
	}

	steps := p.elapsedTime / MicrosecondsToX
	if p.currentStep < steps {
		for _, pipe := range p.Pipes {
			pipe.x -= 1
		}
		p.currentStep = steps
	}

	if len(p.Pipes) > 0 && p.Pipes[0].x < 0 {
		p.Pipes = p.Pipes[1:]
	}
}