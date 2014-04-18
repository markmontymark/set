// see NOTICE for copyright of original material
// I've taken safemap, renamed Insert, Find, Delete to Add, Has, Remove
// and allowed interface{} keys, instead of restricting to string keys

package set_test

import (
    "fmt"
    "github.com/markmontymark/set"
    "sync"
    "testing"
)

func TestSet(t *testing.T) {
    store := set.New()
    fmt.Printf("Initially has %d items\n", store.Len())

    deleted := []int{0, 2, 3, 5, 7, 20, 399, 25, 30, 1000, 91, 97, 98, 99}

    var waiter sync.WaitGroup

    waiter.Add(1)
    go func() { // Concurrent Adder
        for i := 0; i < 100; i++ {
            store.Add(fmt.Sprintf("0x%04X", i))
            if i > 0 && i%15 == 0 {
                fmt.Printf("Added %d items\n", store.Len())
            }
        }
        fmt.Printf("Added %d items\n", store.Len())
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Remover
        for _, i := range deleted {
            key := fmt.Sprintf("0x%04X", i)
            before := store.Len()
            store.Remove(key)
            fmt.Printf("Removed m[%s] (%d) before=%d after=%d\n",
                key, i, before, store.Len())
        }
        waiter.Done()
    }()

    waiter.Add(1)
    go func() { // Concurrent Haser
        for _, i := range deleted {
            for _, j := range []int{i, i + 1} {
                key := fmt.Sprintf("0x%04X", j)
                found := store.Has(key)
                if found {
                    fmt.Printf("Found m[%s] == %d\n", key )
                } else {
                    fmt.Printf("Not found m[%s] (%d)\n", key)
                }
            }
        }
        waiter.Done()
    }()

    waiter.Add(1)
    // can add different types of values into set 
    go func() { // Concurrent Adder
        for i := 0; i < 100; i++ {
            store.Add(i)
            if i > 0 && i%15 == 0 {
                fmt.Printf("Added %d items\n", store.Len())
            }
        }
        fmt.Printf("Added %d items\n", store.Len())
        waiter.Done()
    }()

    waiter.Wait()

    fmt.Printf("Finished with %d items\n", store.Len())
    // not needed here but useful if you want to free up the goroutine
    data := store.Close()
    fmt.Println("Closed")
    fmt.Printf("len == %d\n", len(data))
    //for k, v := range data { fmt.Printf("%s = %v\n", k, v) }
}
