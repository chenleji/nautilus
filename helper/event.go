package helper

import (
	"fmt"
	"github.com/pkg/errors"
	logs "github.com/sirupsen/logrus"
)

const (
	EventPathFmt = "event/%s/%s"
)

var eventMap *EventMap

func GetEventMap() *EventMap {
	if eventMap == nil {
		eventMap = &EventMap{
			data: make(map[string]*Event),
		}
	}
	return eventMap
}

///////////////////////////////////
// event map
//////////////////////////////////

type EventMap struct {
	data map[string]*Event
}

func (m *EventMap) Get(key string) *Event {
	if v, ok := m.data[key]; ok {
		return v
	}
	return nil
}

func (m *EventMap) Registry(e *Event) error {
	if _, ok := m.data[e.Name]; ok {
		return errors.New("已存在同名的event!")
	}
	m.data[e.Name] = e

	consul := Consul{}.New()
	key := fmt.Sprintf(EventPathFmt, Utils{}.GetAppName(), e.Name)

	// start watch
	if ! consul.CheckKey(key, nil) {
		consul.SetKey(key, "")
	}

	e.stopWatchCh = make(chan bool)
	e.stopCallbackCh = make(chan bool)

	e.respCh = consul.WatchKey(key, e.stopWatchCh)

	go func() {
		for {
			select {
			case <-e.stopCallbackCh:
				return

			case resp := <-e.respCh:
				if resp.Error != nil {
					logs.Info("watch event return failed...")
				} else {
					e.Callback(resp.Value)
				}
			}
		}
	}()

	return nil
}

func (m *EventMap) DeRegistry(key string) error {
	if e, ok := m.data[key]; ok {

		// stop watch
		e.stopWatchCh <- true
		e.stopCallbackCh <- true
		delete(m.data, key)
		return nil
	}

	return nil
}

///////////////////////////////////
// event
//////////////////////////////////

type Event struct {
	Name           string
	Callback       func([]byte)
	stopWatchCh    chan bool
	stopCallbackCh chan bool
	respCh         <-chan *WatchResp
}

func (e *Event) Broadcast(value string) error {
	consul := Consul{}.New()
	key := fmt.Sprintf(EventPathFmt, Utils{}.GetAppName(), e.Name)

	if err := consul.SetKey(key, value); err != nil {
		return err
	}

	return nil
}
