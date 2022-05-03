package display

import (
	"go.uber.org/zap"
	"image"
	"image/draw"
	"sync"
)

var logger *zap.SugaredLogger

func NewSubscription(wrappedDisplay Display) *SubscriptionDisplay {
	return &SubscriptionDisplay{
		underlyingDisplay: wrappedDisplay,
		subscribers:       map[string]chan image.Image{},
		lock:              &sync.Mutex{},
	}
}

type SubscriptionDisplay struct {
	underlyingDisplay Display
	subscribers       map[string]chan image.Image
	lock              *sync.Mutex
}

func (v *SubscriptionDisplay) Image() draw.Image {
	return v.underlyingDisplay.Image()
}

func (v *SubscriptionDisplay) Update() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		v.underlyingDisplay.Update()
		wg.Done()
	}()

	v.lock.Lock()
	for id, c := range v.subscribers {
		select {
		case c <- v.Image():
		default:
			logger.Warn("dropping frame for socket", "id", id)
		}
	}
	v.lock.Unlock()

	wg.Wait()
}

func (v *SubscriptionDisplay) RegisterSubscription(id string, imgChan chan image.Image) {
	v.lock.Lock()
	v.subscribers[id] = imgChan
	v.lock.Unlock()
}

func (v *SubscriptionDisplay) RemoveSubscription(id string) {
	v.lock.Lock()
	delete(v.subscribers, id)
	v.lock.Unlock()
}

func (v *SubscriptionDisplay) notify() {

}
