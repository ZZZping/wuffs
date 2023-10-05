// Copyright 2023 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// https://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or https://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build ignore
// +build ignore

package main

// This program exercises the Go standard library's JPEG decoder.
//
// Wuffs' C code doesn't depend on Go per se, but this program gives some
// performance data for specific Go JPEG implementations. The equivalent Wuffs
// benchmarks (on the same test images) are run via:
//
// wuffs bench std/jpeg

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"runtime"
	"time"
)

const (
	iterscale = 20
	reps      = 5
)

type testCase = struct {
	benchname     string
	src           []byte
	itersUnscaled uint32
}

var testCases = []testCase{{
	benchname:     "go_jpeg_decode_19k_8bpp",
	src:           mustLoad("test/data/bricks-gray.jpeg"),
	itersUnscaled: 100,
}, {
	benchname:     "go_jpeg_decode_30k_24bpp_progressive",
	src:           mustLoad("test/data/peacock.progressive.jpeg"),
	itersUnscaled: 50,
}, {
	benchname:     "go_jpeg_decode_30k_24bpp_sequential",
	src:           mustLoad("test/data/peacock.default.jpeg"),
	itersUnscaled: 50,
}, {
	benchname:     "go_jpeg_decode_77k_24bpp",
	src:           mustLoad("test/data/bricks-color.jpeg"),
	itersUnscaled: 30,
}, {
	benchname:     "go_jpeg_decode_552k_24bpp_420",
	src:           mustLoad("test/data/hibiscus.regular.jpeg"),
	itersUnscaled: 5,
}, {
	benchname:     "go_jpeg_decode_552k_24bpp_444",
	src:           mustLoad("test/data/hibiscus.primitive.jpeg"),
	itersUnscaled: 5,
}, {
	benchname:     "go_jpeg_decode_4002k_24bpp",
	src:           mustLoad("test/data/harvesters.jpeg"),
	itersUnscaled: 1,
}}

func mustLoad(filename string) []byte {
	src, err := os.ReadFile("../../" + filename)
	if err != nil {
		panic(err.Error())
	}
	return src
}

func main() {
	if err := main1(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func main1() error {
	fmt.Printf("# Go %s\n", runtime.Version())
	fmt.Printf("#\n")
	fmt.Printf("# The output format, including the \"Benchmark\" prefixes, is compatible with the\n")
	fmt.Printf("# https://godoc.org/golang.org/x/perf/cmd/benchstat tool. To install it, first\n")
	fmt.Printf("# install Go, then run \"go install golang.org/x/perf/cmd/benchstat\".\n")

	for i := -1; i < reps; i++ {
		for _, tc := range testCases {
			runtime.GC()

			start := time.Now()

			iters := uint64(tc.itersUnscaled) * iterscale
			numBytes, err := decode(tc.src)
			if err != nil {
				return err
			}
			for j := uint64(1); j < iters; j++ {
				decode(tc.src)
			}

			elapsedNanos := time.Since(start)

			kbPerS := numBytes * uint64(iters) * 1000000 / uint64(elapsedNanos)

			if i < 0 {
				continue // Warm up rep.
			}

			fmt.Printf("Benchmark%-30s %8d %12d ns/op %8d.%03d MB/s\n",
				tc.benchname, iters, uint64(elapsedNanos)/iters, kbPerS/1000, kbPerS%1000)
		}
	}

	return nil
}

func decode(src []byte) (numBytes uint64, retErr error) {
	m, err := jpeg.Decode(bytes.NewReader(src))
	if err != nil {
		return 0, err
	}

	b := m.Bounds()
	n := uint64(b.Dx()) * uint64(b.Dy())

	// Convert YCbCr to RGBA.
	if _, ok := m.(*image.YCbCr); ok {
		dst := image.NewRGBA(b)
		draw.Draw(dst, b, m, b.Min, draw.Src)
		m = dst
	}

	pix := []byte(nil)
	switch m := m.(type) {
	case *image.Gray:
		n *= 1
	case *image.RGBA:
		n *= 4
		pix = m.Pix
	default:
		return 0, fmt.Errorf("unexpected image type %T", m)
	}

	// Convert RGBA => BGRA.
	if pix != nil {
		for i, iEnd := 0, len(pix)/4; i < iEnd; i += 4 {
			pix[(4*i)+0], pix[(4*i)+2] = pix[(4*i)+2], pix[(4*i)+0]
		}
	}

	return n, nil
}
