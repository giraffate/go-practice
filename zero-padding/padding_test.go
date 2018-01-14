package main

import (
	"fmt"
	"testing"
)

func BenchmarkZeroPadding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = zeroPadding(2018, 8)
	}
}

func BenchmarkFmtZeroPadding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%08d", 2018)
	}
}
