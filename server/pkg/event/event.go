package event

import (
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

var m sync.RWMutex
var eventPool *ants.PoolWithFunc
var handlers map[reflect.Type][]func(i interface{})

func init() {
	var err error
	eventPool, err = ants.NewPoolWithFunc(4, dispatch)
	if err != nil {
		logrus.Error(err)
	}
	handlers = make(map[reflect.Type][]func(i interface{}))
}

func dispatch(i interface{}) {
	handlerList := getHandlerList(i)
	if len(handlerList) == 0 {
		return
	}
	for _, handler := range handlerList {
		handler(i)
	}
}

func Send(e interface{}) {
	if err := eventPool.Invoke(e); err != nil {
		logrus.Error(err)
	}
}

func RegHandler(t reflect.Type, handler func(i interface{})) {
	m.Lock()
	defer m.Unlock()

	handlerList := handlers[t]
	handlerList = append(handlerList, handler)
	handlers[t] = handlerList
}

func getHandlerList(i interface{}) []func(i interface{}) {
	m.RLocker()
	defer m.RUnlock()

	t := reflect.TypeOf(i)
	handlerList, ok := handlers[t]
	if ok {
		return handlerList
	} else {
		logrus.Error("没找到任务处理器，type=" + t.String())
		return nil
	}
}
