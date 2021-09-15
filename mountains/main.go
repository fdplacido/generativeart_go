package main

import (
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

var (
	W = 4000
	H = 2000
	num_mt_ranges = 8
)

type MountainData struct {
	r float64
	g float64
	b float64
	a float64
	w float64
	dc *gg.Context
	big_height_variation float64
	small_height_variation float64
	y_height float64
	big_steps int
	small_steps int
}

func (s *MountainData) InitMountainData() {
	s.y_height = s.big_height_variation * float64(H)
	// set air color
	s.dc.SetRGB(0.7, 0.8, 1)
	s.dc.Clear()

	s.big_steps = 5
	s.small_steps = 5
}

func (s *MountainData) iterate() {
	const factor = 0.9
	s.r = s.r * factor
	s.g = s.g * factor
	s.b = s.b * factor

	// lower height of mountains at every step
	s.y_height = s.y_height + (0.07 * float64(H))
}

func main() {

	rand.Seed(time.Now().Unix())

	mt_range := MountainData{
		r: rand.Float64(),
		g: rand.Float64(),
		b: rand.Float64(),
		// a := rand.Float64() * 0.5 + 0.5,
		a: 1.0,
		w: rand.Float64() * 4 + 1,
		dc: gg.NewContext(W, H),
		big_height_variation: 0.3,
		small_height_variation: 0.1,
	}
	mt_range.InitMountainData()

	for i := 0; i < num_mt_ranges; i++ {
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

	mt.dc.Push()

	// Start from bottom left, to inital left-most point
	mt.dc.LineTo(float64(x_from), float64(dc_h))
	mt.dc.LineTo(float64(x_from), float64(y_from))

	for i := 0; i < mt.big_steps; i++ {

		x_to := x_from + (dc_w/mt.big_steps)
		jitter := float64(y_from) * mt.big_height_variation
		y_to := y_from + (jitter * (-1  + rand.Float64() * 2))

		// small steps should follow the trend of big steps
		x_small_from := x_from
		y_small_from := y_from
		for j := 0; j < mt.small_steps - 1; j++ {
			x_small_to := x_small_from + ((dc_w/mt.big_steps)/mt.small_steps)
			// TODO follow trend
			small_jitter := float64(y_small_from) * mt.small_height_variation
			trend_correction := (y_to - y_from) * 0.3
			y_small_to := y_small_from + (small_jitter * (-1 + rand.Float64() * 2)) + trend_correction
			mt.dc.LineTo(float64(x_small_to), float64(y_small_to))

			// Update small iteration
			x_small_from = x_small_to
			y_small_from = y_small_to
		}

		mt.dc.LineTo(float64(x_to), float64(y_to))

		// Update for next iteration
		x_from = x_to
		y_from = y_to
	}

	// Close polygon going to right-most bottom, then left-most bottom
	mt.dc.LineTo(float64(dc_w), float64(dc_h))
	mt.dc.LineTo(0, float64(dc_h))

	mt.dc.SetFillRule(gg.FillRuleEvenOdd)
	mt.dc.FillPreserve()

	// Paint it
	mt.dc.Stroke()

	mt.dc.Pop()

	drawTree(mt, 500, 1800, 200)
}

func drawTree(mt *MountainData, base_x, base_y, height int) {

	bottom := float64(base_y)
	top := bottom - float64(height)
	base_l := float64(base_x)
	base_r := base_l + (float64(height) * 0.0375)
	half_tree := base_l + ((base_r - base_l) / 2)

	// Draw the trunk
	mt.dc.Push()

	mt.dc.SetRGBA(1, 1, 1, 1)

	mt.dc.LineTo(base_l, bottom)
	mt.dc.LineTo(base_r, bottom)
	mt.dc.LineTo(half_tree, top)
	mt.dc.LineTo(base_l, bottom)

	mt.dc.SetFillRule(gg.FillRuleEvenOdd)
	mt.dc.FillPreserve()
	mt.dc.Stroke()

	mt.dc.Pop()

	// Draw the branches
	branch_base_y := bottom - (float64(height) * 0.125)
	branch_tip_x := float64(height) * 0.25
	branch_side := 1
	branch_reduction_f := branch_tip_x * 0.05

	for branch_base_y > (top + (float64(height) * 0.025)) {
		mt.dc.Push()

		mt.dc.SetRGBA(1, 1, 1, 1)

		mt.dc.LineTo(half_tree, branch_base_y)
		mt.dc.LineTo(half_tree, branch_base_y - (float64(height) * 0.0125))
		mt.dc.LineTo(half_tree + (branch_tip_x * float64(branch_side)) , branch_base_y + (float64(height) * 0.0375))

		// invert branch x position
		branch_side = branch_side * -1
		// reduce branch lenght
		branch_tip_x = branch_tip_x - branch_reduction_f
		// move next branch up
		branch_base_y = branch_base_y - (float64(height) * 0.05)


		mt.dc.SetFillRule(gg.FillRuleEvenOdd)
		mt.dc.FillPreserve()

		// Paint it
		mt.dc.Stroke()

		mt.dc.Pop()
	}
}