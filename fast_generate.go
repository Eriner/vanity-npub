package main

import (
	"encoding/hex"
	"io"
	"math/big"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
	"lukechampine.com/frand"
)

var params = btcec.S256().Params()
var one = new(big.Int).SetInt64(1)
var poolBytes = sync.Pool{New: func() any { return [256/8 + 8]byte{} }}
var poolBigInt = sync.Pool{New: func() any { return new(big.Int) }}

func GenerateFastPrivateKey() string {
	b := poolBytes.Get().([256/8 + 8]byte)
	if _, err := io.ReadFull(frand.Reader, b[:]); err != nil {
		return ""
	}
	k := poolBigInt.Get().(*big.Int).SetBytes(b[:])
	n := poolBigInt.Get().(*big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	poolBytes.Put(b)
	poolBigInt.Put(n)
	defer poolBigInt.Put(k)
	return hex.EncodeToString(k.Bytes())
}
