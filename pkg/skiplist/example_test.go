package skiplist_test

import (
	"fmt"

	"github.com/andremueller/goskiplist/pkg/skiplist"
)

func ExampleNewSkipList() {
	// creates a skip list with key type `int` and value type `string`
	s := skiplist.NewSkipList[int, string]()

	s.Set(1, "cat")
	s.Set(2, "dog")

	x, _ := s.Get(1)
	fmt.Printf("Value: %s", x.Value)
	// Output  "Value: cat"
}

func ExampleSkipList_Set() {
	// creates a skip list with key type `int` and value type `string`
	s := skiplist.NewSkipList[int, string]()

	// sets keys
	s.Set(1, "cat")
	s.Set(2, "dog")

	// overrides the first key
	s.Set(1, "worm")

	x, _ := s.Get(1)
	fmt.Printf("Value: %s", x.Value)
	// Output  "Value: worm"
}
