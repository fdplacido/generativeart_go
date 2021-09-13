package main

import (
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

type MountainData struct {
	r float64
	g float64
	b float64
	a float64
	w float64
	dc *gg.Context
	y_height float64
}

func (s *MountainData) iterate() {
	const factor = 1.2
	s.r = s.r * factor
	s.g = s.g * factor
	s.b = s.b * factor

	// lower height of mountains at every step
	s.y_height = s.y_height + (0.15 * float64(H))
}

var (
	W = 4000
	H = 2000
)

func main() {

	rand.Seed(time.Now().Unix())

	mt_range := MountainData{
		r: rand.Float64() * 0.3,
		g: rand.Float64() * 0.3,
		b: rand.Float64() * 0.3,
		// a := rand.Float64() * 0.5 + 0.5,
		a: 1.0,
		w: rand.Float64() * 4 + 1,
		dc: gg.NewContext(W, H),
		y_height: 0.2 * float64(H),
	}
	mt_range.dc.SetRGB(0, 0, 0)
	mt_range.dc.Clear()

	for i := 0; i < 5; i++ {
		drawMountainLine(&mt_range, W, H)

		mt_range.iterate()
	}

	mt_range.dc.SavePNG("out.png")

}

func drawMountainLine(mt *MountainData, dc_w int, dc_h int) {

	mt.dc.SetRGBA(mt.r, mt.g, mt.b, mt.a)
	mt.dc.SetLineWidth(mt.w)

	x_from := 0
	y_range := 0.1 * float64(dc_h)
	// start at a random place y_range % from starting point
	y_from := float64(mt.y_height) + (y_range * (-1 + rand.Float64() * 2))
	steps := 5

	// Start from bottom left
	mt.dc.LineTo(float64(x_from), float64(dc_h))
	mt.dc.LineTo(float64(x_from), float64(y_from))
	for i := 0; i < steps; i++ {

		x_to := x_from + (dc_w/steps)
		jitter := float64(y_from) * 0.2
		y_to := y_from + (jitter * (-1  + rand.Float64() * 2))

		mt.dc.LineTo(float64(x_to), float64(y_to))

		// Update for next iteration
		x_from = x_to
		y_from = y_to
	}
	mt.dc.LineTo(float64(dc_w), float64(dc_h))
	mt.dc.LineTo(0, float64(dc_h))

	mt.dc.SetFillRule(gg.FillRuleEvenOdd)
	mt.dc.FillPreserve()
	mt.dc.Stroke()

}