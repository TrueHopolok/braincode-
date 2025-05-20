package session

import (
	"crypto/rand"
	"errors"
	"sync"
)

const key_size = 32

type keyChain struct {
	mut sync.Mutex 
	prv []byte
	cur []byte
}

var keys keyChain

func init() {
	keys.mut.Lock()
	defer keys.mut.Unlock()
	keys.prv = make([]byte, key_size)
	keys.cur = make([]byte, key_size)
	rand.Read(keys.prv)
	rand.Read(keys.cur)
}

/*
Updates keys:
 - previous key become previously current key
 - current key randomly generates   
*/
func SwitchKey() {
	keys.mut.Lock() 
	defer keys.mut.Unlock()
	copy(keys.prv, keys.cur)
	rand.Read(keys.cur)
}

// !WARN! - Use only for testing and debuging purposes
/*
Updates keys:
 - previous key become previously current key
 - current key is set to a given key   
*/
func SetKey(new_cur_key []byte) error {
	if len(new_cur_key) != key_size {
		return errors.New("provided key's size is invalid")
	}
	keys.mut.Lock()
	defer keys.mut.Unlock()
	copy(keys.prv, keys.cur)
	copy(keys.cur, new_cur_key)
	return nil
} 
