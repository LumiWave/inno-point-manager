package context

import (
	"sync"
)

const (
	TokenChannel = "TokenChannel"
)

var gContext *ChanContext

var gOnceChan sync.Once

func GetChanInstance() *ChanContext {
	gOnceChan.Do(func() {
		gContext = &ChanContext{}
		context.data = make(map[string]interface{})
	})

	return gContext
}

type ChanContext struct {
	data map[string]interface{}
}

func (o *ChanContext) Put(key string, value interface{}) {
	o.data[key] = value
}

func (o *ChanContext) Get(key string) (interface{}, bool) {
	val, exists := o.data[key]
	return val, exists
}
