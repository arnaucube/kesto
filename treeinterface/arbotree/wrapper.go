// Package arbotree provides the functions for creating and managing an arbo
// merkletree adapted to the CensusTree interface
package arbotree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/arnaucube/arbo"
	"github.com/iden3/go-merkletree/db/leveldb"
	"go.vocdoni.io/dvote/censustree"
)

type Tree struct {
	Tree           *arbo.Tree
	public         uint32
	lastAccessUnix int64 // a unix timestamp, used via sync/atomic
}

func (t *Tree) Init(name, storageDir string) error {
	dbDir := filepath.Join(storageDir, "arbotree.db."+strings.TrimSpace(name))
	storage, err := leveldb.NewLevelDbStorage(dbDir, false) // TODO TMP
	if err != nil {
		return err
	}

	mt, err := arbo.NewTree(storage, 140, arbo.HashFunctionBlake2b) // TODO here the hash function would depend on the usage
	if err != nil {
		return err
	}
	t.Tree = mt
	return nil
}

func (t *Tree) MaxKeySize() int {
	return t.Tree.HashFunction().Len()
}

func (t *Tree) LastAccess() int64 {
	return atomic.LoadInt64(&t.lastAccessUnix)
}

func (t *Tree) updateAccessTime() {
	atomic.StoreInt64(&t.lastAccessUnix, time.Now().Unix())
}

// Publish makes a merkle tree available for queries.
// Application layer should check IsPublish() before considering the Tree available.
func (t *Tree) Publish() {
	atomic.StoreUint32(&t.public, 1)
}

// UnPublish makes a merkle tree not available for queries
func (t *Tree) UnPublish() {
	atomic.StoreUint32(&t.public, 0)
}

// IsPublic returns true if the tree is available
func (t *Tree) IsPublic() bool {
	return atomic.LoadUint32(&t.public) == 1
}

func (t *Tree) Add(index, value []byte) error {
	t.updateAccessTime()
	return t.Tree.Add(index, value)
}

func (t *Tree) AddBatch(indexes, values [][]byte) ([]int, error) {
	t.updateAccessTime()
	return t.Tree.AddBatch(indexes, values)
}

func (t *Tree) GenProof(index, value []byte) ([]byte, error) {
	t.updateAccessTime()
	v, siblings, err := t.Tree.GenProof(index)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(v, value) {
		return nil, fmt.Errorf("value does not match %s!=%s", hex.EncodeToString(v), hex.EncodeToString(value))
	}
	return siblings, nil
}

func (t *Tree) CheckProof(index, value, root, mproof []byte) (bool, error) {
	t.updateAccessTime()
	if root == nil {
		root = t.Root()
	}
	return arbo.CheckProof(t.Tree.HashFunction(), index, value, root, mproof)
}

func (t *Tree) Root() []byte {
	t.updateAccessTime()
	return t.Tree.Root()
}

func (t *Tree) Dump(root []byte) ([]byte, error) {
	return t.Tree.Dump() // TODO pass root once arbo is updated
}

func (t *Tree) DumpPlain(root []byte) ([][]byte, [][]byte, error) {
	t.updateAccessTime()
	var indexes, values [][]byte
	// TODO pass root once arbo is updated
	err := t.Tree.Iterate(func(k, v []byte) {
		if v[0] != arbo.PrefixValueLeaf {
			return
		}
		leafK, leafV := arbo.ReadLeafValue(v)
		indexes = append(indexes, leafK)
		values = append(values, leafV)
	})
	return indexes, values, err
}

func (t *Tree) ImportDump(data []byte) error {
	t.updateAccessTime()
	return t.Tree.ImportDump(data)
}

func (t *Tree) Size(root []byte) (int64, error) {
	count := 0
	err := t.Tree.Iterate(func(k, v []byte) {
		if v[0] != arbo.PrefixValueLeaf {
			return
		}
		count++
	})
	return int64(count), err
}

func (t *Tree) Snapshot(root []byte) (censustree.Tree, error) {
	// TODO
	return t, nil
}

func (t *Tree) HashExists(hash []byte) (bool, error) {
	_, _, err := t.Tree.Get(hash)
	if err != nil {
		return false, err
	}
	return true, nil
}
