package main

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestUnixNano(t *testing.T) {
	ts := time.Now().UnixNano()
	t.Logf("%d", ts)
}

func BenchmarkSystemCallOverhead(b *testing.B) {
	for i := 0; i < 1000; i++ {
		time.Now()
	}

	b.ResetTimer()
	start := time.Now()
	for i := 0; i < b.N; i++ {
		time.Now()
	}
	elapsed := time.Since(start)

	b.ReportMetric(float64(elapsed.Nanoseconds())/float64(b.N), "ns/op")
}

func TestTimeProfile(t *testing.T) {
	f, err := os.Create("time_profile.prof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 5000000; i++ {
		_ = time.Now()
	}
}
