package ratetracker

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Tracker struct {
	Name           string
	ReportInterval time.Duration
	lock           *sync.Mutex
	count          int
	closeChan      chan struct{}
}

func NewTracker(name string) *Tracker {
	t :=
		&Tracker{
			Name:           name,
			ReportInterval: time.Second * 10,
			lock:           &sync.Mutex{},
			count:          0,
			closeChan:      make(chan struct{}),
		}

	go func() {
		for {
			select {
			case <-t.closeChan:
				return
			case <-time.After(t.ReportInterval):
				t.report()
			}
		}
	}()
	return t
}

func (t *Tracker) Close() {
	t.closeChan <- struct{}{}
}

func (t *Tracker) Track() {
	t.lock.Lock()
	t.count += 1
	t.lock.Unlock()
}
func (t *Tracker) report() {
	t.lock.Lock()
	rate := float64(t.count) / t.ReportInterval.Seconds()
	logrus.Infof("%s %f", t.Name, rate)
	t.count = 0
	t.lock.Unlock()
}
