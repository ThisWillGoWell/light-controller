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
	for _, c := range v.subscribers {
		select {
		case c <- v.Image():
		default:
			//zap.S().Warnw("dropping frame for socket", "id", id)
		}
	}
	v.lock.Unlock()

	wg.Wait()
}

func (v *SubscriptionDisplay) RegisterSubscription(id string) chan image.Image {
	imgChan := make(chan image.Image)
	v.lock.Lock()
	v.subscribers[id] = imgChan
	v.lock.Unlock()
	return imgChan
}

func (v *SubscriptionDisplay) RemoveSubscription(id string) {
	v.lock.Lock()
	imageChan := v.subscribers[id]
	delete(v.subscribers, id)
	v.lock.Unlock()
	close(imageChan)
}

func (v *SubscriptionDisplay) notify() {

}
