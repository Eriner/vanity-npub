// generate vanity nostr keys
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"lukechampine.com/frand"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("provide desired prefix as argument")
	}
	desiredPrefix := os.Args[1]
	if strings.ContainsAny(desiredPrefix, "bio1") {
		log.Fatal("prefix cannot contain characters: b, i, o, 1")
	}
	desiredPrefix = "npub1" + desiredPrefix
	var i int
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("\rTries: %d", i)
				<-time.After(1 * time.Second)
			}
		}
	}()
	foundChan := make(chan string, 0)
	bruteforce := func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				i++ // race, but no big deal
				sk := GenerateFastPrivateKey()
				pk, _ := nostr.GetPublicKey(sk)
				npub, _ := nip19.EncodePublicKey(pk)
				if strings.HasPrefix(npub, desiredPrefix) {
					foundChan <- sk
					cancel()
					return
				}
			}
		}
	}
	procs := runtime.GOMAXPROCS(0) - 2
	for i := 0; i < procs; i++ {
		go bruteforce()
	}
	sk := <-foundChan
	<-time.After(100 * time.Millisecond) // prevent stdout stomp
	fmt.Println()
	pk, _ := nostr.GetPublicKey(sk)
	nsec, _ := nip19.EncodePrivateKey(sk)
	npub, _ := nip19.EncodePublicKey(pk)
	fmt.Println("sk:", sk)
	fmt.Println("pk:", pk)
	fmt.Println(nsec)
	fmt.Println(npub)
}

var params = btcec.S256().Params()
var one = new(big.Int).SetInt64(1)
var poolBytes = sync.Pool{New: func() any { return [256/8 + 8]byte{} }}
var poolBigInt = sync.Pool{New: func() any { return new(big.Int) }}

func GenerateFastPrivateKey() string {
	b := poolBytes.Get().([256/8 + 8]byte)
	defer poolBytes.Put(b)
	if _, err := io.ReadFull(frand.Reader, b[:]); err != nil {
		return ""
	}
	//k := new(big.Int).SetBytes(b[:])
	//n := new(big.Int).Sub(params.N, one)
	k := poolBigInt.Get().(*big.Int).SetBytes(b[:])
	n := poolBigInt.Get().(*big.Int).Sub(params.N, one)
	defer poolBigInt.Put(k)
	defer poolBigInt.Put(n)
	k.Mod(k, n)
	k.Add(k, one)
	return hex.EncodeToString(k.Bytes())
}
