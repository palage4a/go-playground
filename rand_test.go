package main_test

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func BenchmarkRandRead(b *testing.B) {
	for i := 0; i < 31; i++ {
		b.Run(ByteCountIEC(1<<i), func(tb *testing.B) {
			read := 0
			for j := 0; j < b.N; j++ {
				buf := make([]byte, 1<<i)
				_, err := rand.Read(buf)
				if err != nil {
					tb.Fatal(err)
				}
				read = read + i<<i
			}
			tb.ReportMetric(float64(read)/float64(tb.Elapsed().Nanoseconds()), "b/ns")
		})
	}
}
