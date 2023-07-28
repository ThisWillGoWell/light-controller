package updater

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"time"
)

var packageLogger = zap.S().With("package", "updater")

type UnmarshalItem func([]byte) (interface{}, error)
type UpdateItem func() (interface{}, error)

type ApiJsonUpdater struct {
	logger         *zap.SugaredLogger
	lastUpdated    time.Time
	fileName       string
	unmarshal      UnmarshalItem
	update         UpdateItem
	data           interface{}
	updateInterval time.Duration
	stop           chan struct{}
	done           chan struct{}
}

func NewUpdater(name string, update UpdateItem, unmarshal UnmarshalItem, interval time.Duration) *ApiJsonUpdater {
	j := &ApiJsonUpdater{
		lastUpdated:    time.Time{},
		fileName:       fmt.Sprintf("%s.json", name),
		unmarshal:      unmarshal,
		update:         update,
		updateInterval: interval,
		stop:           make(chan struct{}),
		logger:         packageLogger.With("name", name),
	}
	j.load()
	j.startUpdate()
	return j
}

func (j *ApiJsonUpdater) Stop() {
	close(j.stop)
	<-j.done
}

func (j *ApiJsonUpdater) load() {
	stats, err := os.Stat(j.fileName)
	if err != nil {
		return
	}
	j.lastUpdated = stats.ModTime()

	f, err := os.Open(j.fileName)
	if err != nil {
		j.logger.Error("failed to open file", zap.Error(err))
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		j.logger.Error("failed to read file", zap.Error(err))
		return
	}

	j.data, err = j.unmarshal(data)
	if err != nil {
		j.logger.Error("failed to unmarshal json file", zap.Error(err))
		return
	}
	return
}

func (j *ApiJsonUpdater) save() {
	data, err := json.MarshalIndent(j.data, "", " ")
	if err != nil {
		j.logger.Error("failed to marshal data", zap.Error(err))
		return
	}
	if err := os.WriteFile(j.fileName, data, 0644); err != nil {
		j.logger.Error("failed to write file", zap.Error(err))
		return
	}
}

func (j *ApiJsonUpdater) startUpdate() {
	update := func() {
		var err error
		j.data, err = j.update()
		if err != nil {
			j.logger.Error("failed to update", zap.Error(err))
		}
		j.save()
	}

	// do we need to update right now
	if j.lastUpdated.Add(j.updateInterval).Before(time.Now()) {
		update()
	}
	go func() {
		// should we wait until we start to update
		if time.Now().Before(j.lastUpdated.Add(j.updateInterval)) {
			<-time.After(time.Until(j.lastUpdated.Add(j.updateInterval)))
		}
		ticker := time.NewTicker(j.updateInterval)
		defer close(j.done)
		for {
			update()
			select {
			case <-ticker.C:
			case <-j.stop:
				return
			}
		}
	}()

}
