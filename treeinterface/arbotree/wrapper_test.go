package arbotree

import (
	"testing"

	"go.vocdoni.io/dvote/censustree"
)

func TestInterface(t *testing.T) {
	storage := t.TempDir()
	tree := &Tree{}
	err := tree.Init("test", storage)
	if err != nil {
		t.Fatal(err)
	}

	var i interface{} = tree
	_, ok := i.(censustree.Tree)
	if !ok {
		t.Fatal("censustree interface not matched by arbotree wrapper")
	}
}

func TestGenProof(t *testing.T) {
	storage := t.TempDir()
	tr1 := &Tree{}
	err := tr1.Init("test1", storage)
	if err != nil {
		t.Fatal(err)
	}

	var keys, values [][]byte
	for i := 0; i < 10; i++ {
		keys = append(keys, []byte{byte(i)})
		values = append(values, []byte{byte(i)})
		err = tr1.Add([]byte{byte(i)}, []byte{byte(i)})
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		p, err := tr1.GenProof(keys[i], values[i])
		if err != nil {
			t.Fatal(err)
		}
		v, err := tr1.CheckProof(keys[i], values[i], tr1.Root(), p)
		if err != nil {
			t.Fatal(err)
		}
		if !v {
			t.Fatal("CheckProof failed")
		}
	}
}
