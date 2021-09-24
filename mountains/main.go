package main

import (
	"math/rand"
	"time"
	// "fmt"

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

type RandomInRange struct {
	min float64
	max float64
}

func (rir *RandomInRange) rnd() float64 {
	return (rir.min + rand.Float64() * (rir.max - rir.min))
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
		drawMountainLine(&mt_range, W, H, i)

		mt_range.iterate()
	}

	mt_range.dc.SavePNG("out.png")

}

func drawMountainLine(mt *MountainData, dc_w int, dc_h int, layer int) {

	sm1 := RandomInRange{min: -1, max: 1}

	mt.dc.SetRGBA(mt.r, mt.g, mt.b, mt.a)
	mt.dc.SetLineWidth(mt.w)

	// Save all intermediate X,Y coordinates to calculate max height of trees later
	horizon_points := make([][2]float64, 0)

	x_from := 0
	y_range := 0.1 * float64(dc_h)
	// start at a random place y_range % from starting point
	y_from := float64(mt.y_height) + (y_range * sm1.rnd())
	horizon_points = append(horizon_points, [2]float64{0: float64(x_from), 1: y_from})

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

			horizon_points = append(horizon_points, [2]float64{0: float64(x_small_from), 1: y_small_from})
		}

		mt.dc.LineTo(float64(x_to), float64(y_to))

		// Update for next iteration
		x_from = x_to
		y_from = y_to

		horizon_points = append(horizon_points, [2]float64{0: float64(x_from), 1: y_from})
	}

	// Close polygon going to right-most bottom, then left-most bottom
	mt.dc.LineTo(float64(dc_w), float64(dc_h))
	mt.dc.LineTo(0, float64(dc_h))

	mt.dc.SetFillRule(gg.FillRuleEvenOdd)
	mt.dc.FillPreserve()

	// Paint it
	mt.dc.Stroke()

	mt.dc.Pop()

	for treeidx := 0; treeidx < 200; treeidx++ {
		// Randomize tree painting
		rndX := RandomInRange{min: 0, max: float64(W)}
		x_num := rndX.rnd()

		// Calculate max Y for tree
		max_y_tree := func() float64 {
			// TODO binary search for X coord
			found_val := 0
			for k := 0; k < len(horizon_points); k++ {
				if x_num > horizon_points[k][0] {
					continue
				} else if x_num <= horizon_points[k][0] {
					found_val = k
					break
				}
			}

			prev_y := horizon_points[found_val-1][1]
			if (prev_y < horizon_points[found_val][1]) {
				return prev_y
			} else {
				return horizon_points[found_val][1]
			}
		}

		rndY := RandomInRange{min: max_y_tree(), max: float64(H)}

		// Calculate height based on layer
		// min 50, max 300
		min_tree_height := 50.0 + (300.0/float64(num_mt_ranges) * float64(layer))
		rndH := RandomInRange{min: min_tree_height, max: min_tree_height + 50}

		drawTree(mt, x_num, rndY.rnd(), rndH.rnd(), layer)
	}

}


type BranchTree struct {
	bottom_start_y float64
	lenght_f float64
	reduction_f float64
	top_limit_y float64
	width_f float64
	inclination_f float64
	separation_f float64
}


func drawTree(mt *MountainData, base_x, base_y, height float64, layer int) {

	sm1 := RandomInRange{min: 0.8, max: 1.2}

	bottom := base_y
	top := bottom - height
	base_l := base_x
	base_r := base_l + ((height * 0.0375) * sm1.rnd())
	half_tree := base_l + ((base_r - base_l) / 2)

	// Draw the trunk
	mt.dc.Push()

	// Make the color whiter based on layer, but only to 0.5
	a := (255.0 - (((255.0/float64(num_mt_ranges))/2) * float64(num_mt_ranges - layer))) / 255.0
	r := (1 - a) * 255 + a * 102
	g := (1 - a) * 255 + a * 51
	b := (1 - a) * 255 + a * 0

	mt.dc.SetRGB255(int(r), int(g), int(b))

	mt.dc.LineTo(base_l, bottom)
	mt.dc.LineTo(base_r, bottom)
	mt.dc.LineTo(half_tree, top)
	mt.dc.LineTo(base_l, bottom)

	mt.dc.SetFillRule(gg.FillRuleEvenOdd)
	mt.dc.FillPreserve()
	mt.dc.Stroke()

	mt.dc.Pop()

	br := BranchTree {
		bottom_start_y: height * 0.125,
		lenght_f: height * 0.25,
		reduction_f: 0.05,
		top_limit_y: height * 0.025,
		width_f: height * 0.0125,
		inclination_f: height * 0.0375,
		separation_f: height * 0.05,
	}

	// Draw the branches
	branch_base_y := bottom - (br.bottom_start_y * sm1.rnd())
	branch_tip_x := br.lenght_f * sm1.rnd()
	branch_side := 1.0
	branch_reduction_f := ((branch_tip_x * br.reduction_f) * sm1.rnd())

	for branch_base_y > (top + br.top_limit_y) {
		mt.dc.Push()

		// Make the color whiter based on layer, but only to 0.5
		a := (255.0 - (((255.0/float64(num_mt_ranges))/2) * float64(num_mt_ranges - layer))) / 255.0
		r := (1 - a) * 255 + a * 0
		g := (1 - a) * 255 + a * 153
		b := (1 - a) * 255 + a * 76
		mt.dc.SetRGB255(int(r), int(g), int(b))

		mt.dc.LineTo(half_tree, branch_base_y)
		mt.dc.LineTo(half_tree, branch_base_y - (br.width_f * sm1.rnd()))
		mt.dc.LineTo(half_tree + (branch_tip_x * branch_side), branch_base_y + (br.inclination_f * sm1.rnd()))
		mt.dc.LineTo(half_tree, branch_base_y)

		// invert branch x position
		branch_side = branch_side * -1
		// reduce branch lenght
		branch_tip_x = branch_tip_x - (branch_reduction_f * sm1.rnd())
		// move next branch up
		branch_base_y = branch_base_y - (br.separation_f * sm1.rnd())


		mt.dc.SetFillRule(gg.FillRuleEvenOdd)
		mt.dc.FillPreserve()

		// Paint it
		mt.dc.Stroke()

		mt.dc.Pop()
	}
}