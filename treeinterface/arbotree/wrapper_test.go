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
