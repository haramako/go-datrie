package trie

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewTrie(t *testing.T) {
	_ = NewTrie()
}

func readTestData() [][]string {
	bin, err := ioutil.ReadFile("test_data.txt")
	if err != nil {
		fmt.Println("Cannot read test_data.txt")
		panic(err)
	}
	txt := string(bin)
	pairs := [][]string{}
	for _, line := range strings.Split(txt, "\n") {
		if line != "" {
			pairs = append(pairs, strings.SplitN(line, " ", 2))
		}
	}
	return pairs
}

func TestStoreAndRetrieve(t *testing.T) {
	trie := NewTrie()

	trie.Store("hoge", 1)

	val, found := trie.Retrieve("hoge")
	if !found {
		t.Errorf("Can't found stored value.")
	}
	if val != 1 {
		t.Errorf("Retrieve value diffred.")
	}
}

func TestStoreData(t *testing.T) {
	data := readTestData()
	trie := NewTrie()

	for i, pair := range data {
		// fmt.Printf("store '%v'\n", pair[0])
		trie.Store(pair[0], int32(i))
	}

	for i, pair := range data {
		num, found := trie.Retrieve(pair[0])
		if !found {
			t.Errorf("Store but not found %v", pair[0])
		}
		if int(num) != i {
			t.Errorf("Invalid result, expect %v but %v with word '%v'", i, num, pair[0])
		}
	}

}

func TestDump(t *testing.T) {
	trie := NewTrie()

	trie.Store("hoge", 1)
	trie.Store("fuga", 2)
	trie.Store("piyo", 3)

	bin := trie.Dump()

	copy_trie, err := Parse(bin)
	if err != nil {
		t.Errorf("cannot load trie from memory: %v", err)
	}

	val, found := copy_trie.Retrieve("hoge")
	if !found {
		t.Errorf("'hoge' not found from loaded trie")
	}
	if val != 1 {
		t.Errorf("'hoge' expect %v but %v", 1, val)
	}

	val, found = copy_trie.Retrieve("fuga")
	if !found {
		t.Errorf("'fuga' not found from loaded trie")
	}
	if val != 2 {
		t.Errorf("'fuga' expect %v but %v", 1, val)
	}

}

func TestWalk(t *testing.T) {
	trie := NewTrie()

	data := []string{"fuga", "hoge", "piyo"}

	for i, str := range data {
		trie.Store(str, int32(i))
	}

	state := trie.Root()

	fmt.Println("!")
	//fmt.Println(state.Current())
	state.Next()
	fmt.Println(state.Current())
	state.Next()
	fmt.Println(state.Current())
	state.Next()
	fmt.Println(state.Current())
}
