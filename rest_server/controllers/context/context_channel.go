package context

import (
	"sync"
)

const (
	Channel_AppPoint = "AppPoint"
)

var gContext *ChanContext

var gOnceChan sync.Once

func GetChanInstance() *ChanContext {
	gOnceChan.Do(func() {
		gContext = &ChanContext{}
		gContext.data = make(map[string]interface{})
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
