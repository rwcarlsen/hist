package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/rwcarlsen/hist"
)

var nbins = flag.Int("nbins", 256, "initial number of bins per dimension")

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("imgsim requires exactly 2 arguments")
	}
	fname1, fname2 := flag.Arg(0), flag.Arg(1)

	f1, err := os.Open(fname1)
	check(err)
	defer f1.Close()

	f2, err := os.Open(fname2)
	check(err)
	defer f2.Close()

	ig1, _, err := image.Decode(f1)
	check(err)

	ig2, _, err := image.Decode(f2)
	check(err)

	img1 := hist.NewDatasetImage(ig1)
	img2 := hist.NewDatasetImage(ig2)

	dist := hist.VarBinDistance(img1, img2, *nbins)
	fmt.Println(dist)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
