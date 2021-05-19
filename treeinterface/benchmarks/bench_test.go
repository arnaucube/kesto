package benchmarks

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/arnaucube/kesto/treeinterface/arbotree"
	asmtree "github.com/p4u/asmt"
	"go.vocdoni.io/dvote/censustree"
	"go.vocdoni.io/dvote/censustree/gravitontree"
)

// func BenchmarkAddBatch(b *testing.B) {
func TestAddBatch(t *testing.T) {
	nLeafs := 100_000
	fmt.Printf("nCPU: %d, nLeafs: %d\n", runtime.NumCPU(), nLeafs)

	// prepare inputs
	var ks, vs [][]byte
	for i := 0; i < nLeafs; i++ {
		k := randomBytes(32)
		v := randomBytes(32)
		ks = append(ks, k)
		vs = append(vs, v)
	}

	tree1 := &asmtree.Tree{}
	benchmarkAdd(t, "asmtree", tree1, ks, vs)

	tree1 = &asmtree.Tree{}
	benchmarkAddBatch(t, "asmtree", tree1, ks, vs)

	tree2 := &gravitontree.Tree{}
	benchmarkAdd(t, "gravitontree", tree2, ks, vs)

	tree2 = &gravitontree.Tree{}
	benchmarkAddBatch(t, "gravitontree", tree2, ks, vs)

	tree3 := &arbotree.Tree{}
	benchmarkAdd(t, "arbo", tree3, ks, vs)

	tree3 = &arbotree.Tree{}
	benchmarkAddBatch(t, "arbo", tree3, ks, vs)
}

func benchmarkAdd(t *testing.T, name string, tree censustree.Tree, ks, vs [][]byte) {
	storage := t.TempDir()
	err := tree.Init("test1", storage)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	for i := 0; i < len(ks); i++ {
		if err := tree.Add(ks[i], vs[i]); err != nil {
			t.Fatal(err)
		}
	}
	printRes(t, name+".Add loop", time.Since(start))
}

func benchmarkAddBatch(t *testing.T, name string, tree censustree.Tree, ks, vs [][]byte) {
	storage := t.TempDir()
	err := tree.Init("test1", storage)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	invalids, err := tree.AddBatch(ks, vs)
	if err != nil {
		t.Fatal(err)
	}
	if len(invalids) != 0 {
		t.Fatal("len(invalids)!=0")
	}
	printRes(t, name+".AddBatch", time.Since(start))
}

func printRes(t *testing.T, name string, duration time.Duration) {
	fmt.Printf("	%s:	%s \n", name, duration)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
