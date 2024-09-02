package envs

import (
	"testing"
)

func BenchmarkLoadEnv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Load()
	}
}
