package main

import (
	"testing"

	"github.com/nbd-wtf/go-nostr"
)

func BenchmarkGeneratePrivateKey(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nostr.GeneratePrivateKey()
	}
}

func BenchmarkGenerateFastPrivateKey(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateFastPrivateKey()
	}
}
