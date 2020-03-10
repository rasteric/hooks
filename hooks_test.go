package hooks

import (
	"math/rand"
	"testing"
	"time"
)

var x int

func TestHooks(t *testing.T) {
	Add(1, func(a []interface{}) {
		x = 2
	})
	id2 := Add(1, func(a []interface{}) {
		x = 3
	})
	id3 := Add(2, func(a []interface{}) {
		x = 4
	})
	Exec(1)
	if x != 2 && x != 3 {
		t.Errorf("Exec(1): expected x to have value 2 or 3 after Add, given %v", x)
	}
	Exec(2)
	if x != 4 {
		t.Errorf("Exec(2): expected x to have value 4 after Add, given %v", x)
	}
	Remove(2, id3)
	Exec(1)
	Exec(2)
	if x != 2 && x != 3 {
		t.Errorf("Exec(2): expected x to have value 2 or 3 after Remove(2), given %v", x)
	}
	RemoveAll(2)
	Remove(1, id2)
	Exec(1)
	if x != 2 {
		t.Errorf("Exec(1): expected x to have value 2, given %v", x)
	}
	Add(1, func(a []interface{}) {
		x = 3
	})
	for i := 0; i < 20000; i++ {
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
			Exec(1)
			Exec(2)
			Exec(3)
		}()
	}

}
