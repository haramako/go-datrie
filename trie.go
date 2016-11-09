package trie

import (
	//#cgo CFLAGS: -I ..
	//#include "datrie/trie.h"
	//#include "datrie/alpha-map.h"
	"C"
	"bytes"
	"fmt"
	"runtime"
	"unsafe"
)

type Trie struct {
	trie *C.Trie
}

type TrieState struct {
	trie  *Trie
	state *C.TrieState
}

type TrieIterator struct {
	state    *TrieState
	iterator *C.TrieIterator
}

func finalizeTrie(t *Trie) {
	if t.trie != nil {
		C.trie_free(t.trie)
	}
}

func NewTrie() *Trie {
	t := &Trie{}
	runtime.SetFinalizer(t, finalizeTrie)
	alpha_map := C.alpha_map_new()
	C.alpha_map_add_range(alpha_map, 1, 255)
	t.trie = C.trie_new(alpha_map)
	return t
}

func (t *Trie) Dump() []byte {
	buf := make([]byte, 100000)
	result := C.trie_write_to_memory(t.trie, unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	if result < 0 {
		panic("error!")
	} else {
		return buf[0:int(result)]
	}
}

func Parse(buf []byte) (*Trie, error) {
	ctrie := C.trie_new_from_memory(unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	if ctrie == nil {
		return nil, fmt.Errorf("cannot load trie")
	}
	t := &Trie{trie: ctrie}
	runtime.SetFinalizer(t, finalizeTrie)
	return t, nil
}

func (t *Trie) Store(key string, val int32) error {
	runes := append(bytes.Runes([]byte(key)), 0)
	buf := (*C.AlphaChar)(unsafe.Pointer(&runes[0]))
	if C.trie_store(t.trie, buf, C.TrieData(val)) == 0 {
		return fmt.Errorf("cannot store %v", key)
	}
	return nil
}

func (t *Trie) Retrieve(key string) (int32, bool) {
	runes := append(bytes.Runes([]byte(key)), 0)
	buf := (*C.AlphaChar)(unsafe.Pointer(&runes[0]))
	var result C.TrieData
	if C.trie_retrieve(t.trie, buf, &result) == 0 {
		return 0, false
	}
	return int32(result), true
}

func (t *Trie) Root() *TrieIterator {
	return newTrieIterator(newTrieState(t))
}

func newTrieState(trie *Trie) *TrieState {
	s := &TrieState{trie: trie, state: C.trie_root(trie.trie)}
	runtime.SetFinalizer(s, finalizeTrieState)
	return s
}

func finalizeTrieState(s *TrieState) {
	if s.state != nil {
		C.trie_state_free(s.state)
	}
}

func newTrieIterator(state *TrieState) *TrieIterator {
	citerator := C.trie_iterator_new(state.state)
	it := &TrieIterator{state: state, iterator: citerator}
	runtime.SetFinalizer(it, finalizeTrieIterator)
	return it
}

func finalizeTrieIterator(it *TrieIterator) {
	if it.iterator != nil {
		C.trie_iterator_free(it.iterator)
	}
}

func (it *TrieIterator) Next() bool {
	return C.trie_iterator_next(it.iterator) != 0
}

func (it *TrieIterator) Current() (string, int) {
	chars := C.trie_iterator_get_key(it.iterator)
	if chars == nil {
		return "", 0
	}
	data := C.trie_iterator_get_data(it.iterator)
	return chars.String(), int(data)
}

func (c *C.AlphaChar) String() string {
	runes := *(*[256]rune)(unsafe.Pointer(c))
	var i int
	for i = 0; runes[i] != 0; i++ {
	}
	return string(runes[:i])
}
