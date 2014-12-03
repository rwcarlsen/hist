package hist

import (
	"image"
	"math"
)

type Dataset interface {
	// At returns a series of values for point i.  The values must be in the
	// range [0, 1).
	At(i int) []float64
	// Dims returns the number of values associated with each point in the
	// dataset.
	Dims() int
	// Len returns the number of points in the dataset.
	Len() int
}

// VarBinDistance returns variable bin with distance between d1 and d2 with
// the given initial number of bins.
func VarBinDistance(d1, d2 Dataset, nbins int) float64 {
	if d1.Dims() != d1.Dims() {
		panic("datasets don't have same number of dimensions")
	}

	bounds := make([]Bins, d1.Dims())
	for i := range bounds {
		bounds[i] = NewBins(0, 1, nbins)
	}
	h1 := Make(d1, bounds)
	h2 := Make(d2, bounds)

	distances := []float64{}

	for nbins > 0 {
		d := L1Distance(h1, h2)
		distances = append(distances, d)

		nbins /= 2
		inter := Intersect(h1, h2)
		h1 = Diff(h1, inter)
		h2 = Diff(h2, inter)
		bounds := make([]Bins, d1.Dims())
		for i := range bounds {
			bounds[i] = NewBins(0, 1, nbins)
		}
	}

	tot := 0.0
	for _, d := range distances {
		tot += d
	}
	return tot / float64(len(distances))
}

func Diff(h1, h2 Hist) Hist {
	diff := Hist{}
	for k := range h1 {
		diff[k] = h1[k] - h2[k]
	}
	for k := range h2 {
		diff[k] = h1[k] - h2[k]
	}
	return diff
}

func Intersect(h1, h2 Hist) Hist {
	inter := Hist{}
	for k := range h1 {
		inter[k] = math.Min(h1[k], h2[k])
	}
	for k := range h2 {
		inter[k] = math.Min(h1[k], h2[k])
	}
	return inter
}

func L2Distance(h1, h2 Hist) float64 {
	distances := Hist{}
	for k := range h1 {
		distances[k] = h1[k] - h2[k]
	}
	for k := range h2 {
		distances[k] = h1[k] - h2[k]
	}
	d := 0.0
	for _, val := range distances {
		d += val * val
	}
	return d
}

func L1Distance(h1, h2 Hist) float64 {
	distances := Hist{}
	for k := range h1 {
		distances[k] = h1[k] - h2[k]
	}
	for k := range h2 {
		distances[k] = h1[k] - h2[k]
	}
	d := 0.0
	for _, val := range distances {
		d += math.Abs(val)
	}
	return d
}

type Hist map[uint64]float64

func (h Hist) Freq(bounds []Bins, pos ...int) float64 {
	return h[Key(bounds, pos...)]
}

// Bins holds a set of bounds for the bins of a particular dimension.  The
// length of bins should be 1 larger than the number of bins.
type Bins []float64

func NewBins(start, end float64, nbins int) Bins {
	v := start
	w := (end - start) / float64(nbins)
	bins := Bins{}
	for i := 0; i <= nbins; i++ {
		bins = append(bins, v)
		v += w
	}
	bins[len(bins)-1] = end
	return bins
}

// Bin returns the index of the bin that val falls into using bounds specified
// in bins.
func (bs Bins) Bin(val float64) int {
	for b := range bs[:len(bs)-1] {
		if bs[b] <= val && val < bs[b+1] {
			return b
		}
	}
	return len(bs) - 1
}

// Make generates a normalized frequency histogram of the multi-dimensional data
// in d placed into multi-dimensional bins defined by bounds.
func Make(d Dataset, bounds []Bins) Hist {
	if d.Dims() != len(bounds) {
		panic("datasets don't have same number of dimensions as bin bounds")
	}

	h := Hist{}
	for i := 0; i < d.Len(); i++ {
		vals := d.At(i)
		pos := make([]int, d.Dims())
		for j, val := range vals {
			pos[j] = bounds[j].Bin(val)
		}
		h[Key(bounds, pos...)] += 1
	}

	for k := range h {
		h[k] /= float64(d.Len())
	}
	return h
}

// Key returns a unique uint64 key for a specific bin over the binset domain.
// The product of all the number of bins in the binset must be less than
// MaxUint64.
func Key(bounds []Bins, pos ...int) uint64 {
	key := uint64(0)
	mult := uint64(1)
	for _, i := range pos {
		key += uint64(i) * mult
		mult *= uint64(i)
	}
	return key
}

type Image struct {
	image.Image
	w, h   int
	x0, y0 int
}

func NewDatasetImage(img image.Image) *Image {
	r := img.Bounds()
	return &Image{
		Image: img,
		w:     r.Max.X - r.Min.X,
		h:     r.Max.Y - r.Min.Y,
		x0:    r.Min.X,
		y0:    r.Min.Y,
	}
}

func (img *Image) Len() int {
	return img.w * img.h
}

func (img *Image) Dims() int { return 5 }

func (img Image) At(i int) []float64 {
	dy := i / img.w
	dx := i % img.w
	y := dy + img.y0
	x := dx + img.x0

	r, g, b, _ := img.Image.At(x, y).RGBA()

	return []float64{
		float64(dx) / float64(img.w),
		float64(dy) / float64(img.h),
		float64(r) / 0xFFFF,
		float64(g) / 0xFFFF,
		float64(b) / 0xFFFF,
	}
}
