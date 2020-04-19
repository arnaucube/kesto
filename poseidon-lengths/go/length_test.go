package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

func TestLen(t *testing.T) {
	m := []byte("45")
	h, err := poseidon.HashBytes(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("bigint", h.String())
	fmt.Println("length", len(h.Bytes()))
	fmt.Println("bytes", h.Bytes())
	fmt.Println("hex", hex.EncodeToString(h.Bytes()))
	if len(h.Bytes()) != 31 {
		t.Fatal("expected 31 bytes")
	}

}
