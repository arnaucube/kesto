package main

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/arnaucube/go-blindsecp256k1"
	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func calculateKeysFromK(t *testing.T, k *big.Int) (*btcec.PublicKey, *blindsecp256k1.PublicKey) {
	curve := btcec.S256()

	x, y := curve.ScalarBaseMult(k.Bytes())

	skEth := (*btcec.PrivateKey)(&ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: k,
	})
	pkEthAux := (*btcec.PublicKey)(&skEth.PublicKey)
	pkEth := skEth.PubKey()
	assert.Equal(t, pkEth, pkEthAux)

	skBlind := blindsecp256k1.PrivateKey(*k)
	pkBlind := skBlind.Public()

	return pkEth, pkBlind
}

func TestCompatibleKeys(t *testing.T) {
	// using incremental k
	for i := 0; i < 10000; i++ {
		k := big.NewInt(int64(i))

		pkEth, pkBlind := calculateKeysFromK(t, k)
		assert.Equal(t, pkEth.X.String(), pkBlind.X.String())
		assert.Equal(t, pkEth.Y.String(), pkBlind.Y.String())
		// fmt.Println("k", k)
		// fmt.Println(pkEth)
		// fmt.Println(pkBlind)
	}

	// using random k
	for i := 0; i < 10000; i++ {
		kaux, err := btcec.NewPrivateKey(btcec.S256())
		require.NoError(t, err)
		k := kaux.D

		pkEth, pkBlind := calculateKeysFromK(t, k)
		assert.Equal(t, pkEth.X.String(), pkBlind.X.String())
		assert.Equal(t, pkEth.Y.String(), pkBlind.Y.String())
		// fmt.Println("k", k)
		// fmt.Println(pkEth)
		// fmt.Println(pkBlind)
	}
}
