package hist

import "image"

type Dataset interface {
	At(i int) []float64
	Dims() int
	Len() int
}

func VarBinDistance(d1, d2 Dataset, bins [][]float64) float64 {
	if d1.Dims() != d1.Dims() {
		panic("datasets don't have same number of dimensions")
	} else if d1.Dims() != len(bins) {
		panic("datasets don't have same number of dimensions as bin sets")
	}

	h := &Hist{Bins: bins, Freqs: map[uint64]float64{}}
	for i := 0; i < d1.Len(); i++ {
		vals := d1.At(i)
		pos := make([]int, d1.Dims())
		for j, val := range vals {
			pos[j] = Bin(val, bins[j])
		}
		h.Freqs[Key(bins, pos...)] += 1
	}
}

type Hist map[uint64]float64

func (h Hist) Freq(pos ...int) float64 {
	return h[Key(pos...)]
}

// Bins holds a set of bounds for the bins of a particular dimension.  The
// length of bins should be 1 larger than the number of bins.
type Bins []float64

// Freq returns a normalized frequency histogram of the multi-dimensional data
// in d placed into multi-dimensional bins defined by bounds.
func Freq(d Dataset, bounds []Bins) *Hist {
	if d.Dims() != len(bounds) {
		panic("datasets don't have same number of dimensions as bin bounds")
	}

	h := &Hist{}
	for i := 0; i < d.Len(); i++ {
		vals := d.At(i)
		pos := make([]int, d.Dims())
		for j, val := range vals {
			pos[j] = Bin(val, bounds[j])
		}
		h[Key(bounds, pos...)] += 1
	}

	for k, val := range h {
		h[k] /= float64(d.Len())
	}
	return h
}

// Bin returns the index of the bin that val falls into using bounds specified
// in bins.
func Bin(val float64, bins Bins) int {
	for b, bound := range bins[:len(bins)-1] {
		if bins[b] <= val && val < bins[b+1] {
			return b
		}
	}
	panic("val is out of bins range")
}

// Key returns a unique uint64 key for a specific bin over the binset domain.
// The product of all the number of bins in the binset must be less than
// MaxUint64.
func Key(bounds []Bins, pos ...int) uint64 {
	key := 0
	mult := 1
	for _, i := range pos {
		key += i * mult
		mult *= i
	}
	return key
}

type Image struct {
	image.Image
	w, h   int
	x0, y0 int
}

func NewImage(img image.Image) *Image {
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
	return w * h
}

func (img *Image) Dims() int { return 5 }

func (img Image) At(i int) []float64 {
	dy := i / w
	dx := i % w
	y := dy + img.y0
	x := dx + img.x0

	r, g, b, _ := img.Image.At(x, y).RGBA()

	return []float64{
		float64(dx) / float64(w),
		float64(dy) / float64(h),
		float64(r) / 0xFFFF,
		float64(g) / 0xFFFF,
		float64(b) / 0xFFFF,
	}
}
