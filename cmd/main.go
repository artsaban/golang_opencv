package main

import (
	"fmt"
	"image"
	"os"
	"sort"

	colorful "github.com/lucasb-eyer/go-colorful"
	cv "gocv.io/x/gocv"
)

type ColorFrequency map[string]uint

func rankColorsByFrequency(colorFrequencies ColorFrequency) PairList {
	pl := make(PairList, len(colorFrequencies))
	i := 0
	for k, v := range colorFrequencies {
		pl[i] = Pair{k, v, false}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// TODO: это уже не пара получается :harold: Нужно переименовать
type Pair struct {
	Key        string
	Value      uint
	Compressed bool
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	imgPath := os.Args[1]

	img := cv.IMRead(imgPath, cv.IMReadColor)
	cv.Resize(img, &img, image.Point{X: 640, Y: 480}, 0, 0, cv.InterpolationArea)
	cv.CvtColor(img, &img, cv.ColorBGRToRGB)

	colorFrequencies, colorRgbs := createColorMappings(&img)
	rankedColors := rankColorsByFrequency(colorFrequencies)
	compressed := compressColors(rankedColors, colorRgbs, 12)

	fmt.Println(compressed)
}

// TODO: из этой функции возвращать PairList (но в Pair докинуть еще данных -- цвет в виде colorful.Color)
func createColorMappings(img *cv.Mat) (ColorFrequency, map[string]colorful.Color) {
	counts := make(ColorFrequency)
	colors := make(map[string]colorful.Color)

	size := img.Size()
	height, width := size[0], size[1]
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			pixelColor := img.GetVecbAt(x, y)
			red := pixelColor[0]
			green := pixelColor[1]
			blue := pixelColor[2]
			colorKey := fmt.Sprintf("%d%d%d", red, green, blue)
			if _, ok := counts[colorKey]; !ok {
				counts[colorKey] = 1
				colors[colorKey] = colorful.Color{R: float64(red), G: float64(green), B: float64(blue)}
			} else {
				counts[colorKey] += 1
			}
		}
	}

	return counts, colors
}

func compressColors(colors PairList, clrs map[string]colorful.Color, tolerance float64) PairList {
	if tolerance <= 0 {
		return colors
	}

	i := 0
	for i < len(colors) {
		larger := &colors[i]
		if !larger.Compressed {
			j := i + 1
			for j < len(colors) {
				smaller := &colors[j]
				if !smaller.Compressed && clrs[smaller.Key].DistanceCIE76(clrs[larger.Key]) <= tolerance {
					larger.Value += smaller.Value
					smaller.Compressed = true
				}
				j += 1
			}
		}
		i += 1
	}

	pl := make(PairList, 0)
	for _, color := range colors {
		if !color.Compressed {
			pl = append(pl, color)
		}
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
