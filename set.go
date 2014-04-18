// set stores values as interface{} instead of providing specific type variants like 
// stringset, intset, <YourSuperComplexDataType>set, because user code can just as 
// cast the output of Close() from map[interface{}]struct{} to map[YourSuperComplexDataType]struct{}
// very easily for the few the cases where you need that.  It's just as easy to get, set, print interface{}
// values without having to do any casting 
package set

type safeSet chan commandData

type commandData struct {
    action  commandAction
    key     interface{}
    result  chan<- interface{}
    data    chan<- map[interface{}]struct{}
}

type commandAction int

const (
    remove commandAction = iota
    end
    has
    add
    length
)

type Interface interface {
    Add(interface{})
    Remove(interface{})
    Has(interface{}) bool
    Len() int
    Close() map[interface{}]struct{}
}

func New() Interface {
    sm := make(safeSet) // type safeSet chan commandData
    go sm.run()
    return sm
}

func (sm safeSet) run() {
    store := make(map[interface{}]struct{})
    for command := range sm {
        switch command.action {
			  case add:
					store[command.key] = struct{}{}
			  case remove:
					delete(store, command.key)
			  case has:
					_, found := store[command.key]
					command.result <- found
			  case length:
					command.result <- len(store)
			  case end:
					close(sm)
					command.data <- store
		  }
	 }
}

func (sm safeSet) Add(key interface{}) {
    sm <- commandData{action: add, key: key}
}

func (sm safeSet) Remove(key interface{}) {
    sm <- commandData{action: remove, key: key}
}

func (sm safeSet) Has(key interface{}) bool {
    reply := make(chan interface{})
    sm <- commandData{action: has, key: key, result: reply}
    return (<-reply).(bool)
}

func (sm safeSet) Len() int {
    reply := make(chan interface{})
    sm <- commandData{action: length, result: reply}
    return (<-reply).(int)
}

// Close() may only be called once per safe map; all other methods can be
// called as often as desired from any number of goroutines
func (sm safeSet) Close() map[interface{}]struct{} {
    reply := make(chan map[interface{}]struct{})
    sm <- commandData{action: end, data: reply}
    return <-reply
}
