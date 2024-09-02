package main

import "testing"

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initLogging()
	}
}
