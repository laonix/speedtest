package speedtest

import "testing"

// benchmarks

func BenchmarkDefaultRunner_Run_ookla(b *testing.B) {
	r := NewRunner()

	for i := 0; i < b.N; i++ {
		_, _ = r.Run("ookla")
	}
}

func BenchmarkDefaultRunner_Run_netflix(b *testing.B) {
	r := NewRunner()

	for i := 0; i < b.N; i++ {
		_, _ = r.Run("netflix")
	}
}
