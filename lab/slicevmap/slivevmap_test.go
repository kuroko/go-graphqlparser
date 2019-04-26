package slicevmap

import (
	"fmt"
	"testing"
)

type brt struct {
	m map[int]struct{}
	s []int
}

func BenchmarkMap(b *testing.B) {
	var br brt

	ranges := []int{1e1, 1e2, 1e3, 1e4, 1e5}

	for _, r := range ranges {
		b.Run(fmt.Sprintf("insert-%d", r), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				br.m = make(map[int]struct{})
				for j := 0; j < r; j++ {
					br.m[j] = struct{}{}
					_ = br.m
				}
			}
		})

		b.Run(fmt.Sprintf("access-%d", r), func(b *testing.B) {
			br.m = make(map[int]struct{})
			for j := 0; j < r; j++ {
				br.m[j] = struct{}{}
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				v := br.m[r-1]
				_ = v
			}
		})
	}
}

func BenchmarkSlice(b *testing.B) {
	var br brt

	ranges := []int{1e1, 30, 50, 1e2, 1e3, 1e4, 1e5}

	for _, r := range ranges {
		b.Run(fmt.Sprintf("insert-%d", r), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				br.s = []int{}
				for j := 0; j < r; j++ {
					br.s = append(br.s, j)
				}
			}
		})

		b.Run(fmt.Sprintf("access-%d", r), func(b *testing.B) {
			br.s = []int{}
			for j := 0; j < r; j++ {
				br.s = append(br.s, j)
			}

			b.ReportAllocs()
			b.ResetTimer()

			// Worst case scenario for slice.
			for i := 0; i < b.N; i++ {
				for _, v := range br.s {
					if v == (r - 1) {
						break
					}
				}
			}
		})
	}
}
