# Vanity npub

A vanity prefix bech-32 nostr generator.

## What

Do you want a nostr public key like: 

npub1**swag**7ws4yul24p0pf7cacn7cghqkutdnm35z075vy68ggqpqjcyswn8ekc?

If so, this is the tool for you.

## How

```
$ go install github.com/eriner/vanity-npub

$ vanity_npub swag
```

## Notes

The characters `b`, `i`, `o`, and `1` are not valid bech-32 characters, so they may not be used in the vanity prefix.

## Benchmarks

This program uses a modified version of the [nbd-wtf/go-nostr](https://github.com/nbd-wtf/go-nostr)'s `GeneratePrivateKey` function, intended for speed.

Namely, `sync.Pool`s are used to reduce memory allocations and [frand](https://github.com/lukechampine/frand) (a CSPRNG) is used instead of `crypto/rand`.

The benchmarks below show an approximate 2x speed increase while maintaining the cryptographic security of key generation.

```
% go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/eriner/vanity-npub
cpu: AMD Ryzen 7 2700X Eight-Core Processor         
BenchmarkGeneratePrivateKey-16        	 499538	     2073 ns/op
BenchmarkGenerateFastPrivateKey-16    	1085912	     1146 ns/op
PASS
ok  	github.com/eriner/vanity-npub	3.240s
```
