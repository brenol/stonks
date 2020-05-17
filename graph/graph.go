package graph

import (
	"github.com/ericm/stonks/api"
	"github.com/shopspring/decimal"
)

func borderHorizontal(out *string, width int) {
	for _i := 0; _i < width-2; _i++ {
		*out += "━"
	}
}

// GenerateGraph with ASCII graph with ANSI escapes
func GenerateGraph(chart *api.Chart, width int, height int) (string, error) {
	out := "┏"
	borderHorizontal(&out, width)
	out += "┓"
	// interval, err := time.ParseDuration(string(chart.Interval))
	// if err != nil {
	// 	return "", err
	// }
	matrix := make([][]*api.Bar, height)
	for i := range matrix {
		matrix[i] = make([]*api.Bar, width)
	}
	ran := chart.High.Sub(chart.Low)
	spacing := (width) / (chart.Length)
	out += "\n"
	var last *api.Bar
	for x, bar := range chart.Bars {
		bar.Char = "─"
		y := int(bar.Current.Sub(chart.Low).Div(ran).Mul(
			decimal.NewFromInt((int64(height)))).Floor().IntPart())
		matrix[y][x*spacing] = bar
		bar.Y = y
		if last != nil {
			next := last.Y - bar.Y
			var char string
			currY := last.Y
			switch {
			case next > 0:
				char = "╱"
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY--
					if currY >= 0 && currY >= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next < 0:
				char = "╲"
				bar.Char = char
				for i := 0; i < spacing-1; i++ {
					currY++
					if currY < height && currY <= y {
						matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
					}
				}
			case next == 0:
				char = "─"
				last.Char = char
				for i := 0; i < spacing-1; i++ {
					matrix[currY][i+((x-1)*spacing)+1] = &api.Bar{Char: char}
				}
			}
			// Edge cases
			switch last.Char {
			case "╱":
				switch char {
				case "╲":
					if matrix[y][(x*spacing)-1] != nil {
						last.Char = "▁"
					} else {
						last.Char = "ʌ"
					}
				case "╱":
					last.Char = "╱"
				}
			case "╲":
				switch char {
				case "╲":
					if matrix[y][(x*spacing)-1] != nil {
						last.Char = "▔"
					} else {
						last.Char = "▁"
					}
				case "╱":
					last.Char = "╱"
				}
			}
		}
		last = bar
	}
	for _, slc := range matrix {
		out += "┃"
		for _, ptr := range slc {
			if ptr != nil {
				out += ptr.Char
			} else {
				out += " "
			}
		}
		out += "\n"
	}
	return out, nil
}
