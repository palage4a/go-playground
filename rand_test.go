package main_test

import (
	"crypto/rand"
	"testing"
)

func BenchmarkRandRead(b *testing.B) {
	for _, tc := range []struct {
		name string
		size int
	}{
		{"1b", 1},
		{"512b", 1 << 9},
		{"1kb", 1 << 10},
		{"2kb", 1 << 11},
		{"512kb", 1 << 19},
		{"1mb", 1 << 20},
		{"2mb", 1 << 21},
		{"512mb", 1 << 29},
		{"1gb", 1 << 30},
	} {

		b.Run(tc.name, func(tb *testing.B) {
			read := 0
			for i := 0; i < b.N; i++ {
				buf := make([]byte, tc.size)
				_, err := rand.Read(buf)
				if err != nil {
					tb.Fatal(err)
				}
				read += tc.size

			}
			tb.ReportMetric(float64(read)/float64(tb.Elapsed().Nanoseconds()), "b/ns")
		})
	}
}
