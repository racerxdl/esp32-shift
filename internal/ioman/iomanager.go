package ioman

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/quan-to/slog"
	"github.com/racerxdl/esp32-shift/pkg/proto"
	"github.com/racerxdl/esp32-shift/pkg/queue"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var log = slog.Scope("IO")

type iomanager struct {
	queue     queue.Manager
	dev       shift.ESPShift
	running   bool
	devLock   sync.Mutex
	stateLock sync.Mutex
	baseTopic string

	pinStates []shift.IOLevel
	status    shift.BoardStatus
}

func MakeIOManager(dev shift.ESPShift, baseTopic string, q queue.Manager) io.Closer {
	iom := &iomanager{
		queue:     q,
		dev:       dev,
		running:   true,
		devLock:   sync.Mutex{},
		stateLock: sync.Mutex{},
		baseTopic: baseTopic,
	}

	q.SetOnMessage(iom.messageHandle)

	topic := fmt.Sprintf("%s/#", baseTopic)
	err := q.Subscribe(topic)
	if err != nil {
		log.Error("Cannot subscribe to topic %q: %s", topic, err)
	}

	go iom.loop()

	return iom
}

func (iom *iomanager) resizePinState(size int) {
	iom.stateLock.Lock()
	size += 1
	if len(iom.pinStates) < size {
		newStates := make([]shift.IOLevel, size)
		copy(newStates, iom.pinStates)
		iom.pinStates = newStates
	}
	iom.stateLock.Unlock()
}

func (iom *iomanager) messageHandle(msg mqtt.Message) {
	topic := msg.Topic()
	value := string(msg.Payload())
	//log.Debug("Received Message: %s => %s", topic, value)
	t := strings.Split(topic, "/")
	if len(t) <= 1 {
		return
	}

	level := shift.LOW
	if value == "1" || value == "true" {
		level = shift.HIGH
	}

	pinNum, err := strconv.ParseInt(t[1], 10, 32)

	if err != nil {
		log.Error("Error parsing topic pin %q: %s", t[1], err)
		return
	}

	iom.resizePinState(int(pinNum))

	if iom.pinStates[int(pinNum)] != level {
		iom.devLock.Lock()
		err = iom.dev.SetPin(uint8(pinNum), level)
		iom.devLock.Unlock()
		iom.pinStates[int(pinNum)] = level
	}

	if err != nil {
		log.Error("Error setting pin %d to %s: %s", pinNum, level, err)
	}
}

func (iom *iomanager) updateStatus() {
	s := make([]int, len(iom.status))
	for i, v := range iom.status {
		s[i] = int(v)
	}
	topic := fmt.Sprintf("%s_BS", iom.baseTopic)
	d, _ := json.Marshal(s)
	//log.Debug("Publishing status to %q: %s", topic, string(d))
	_ = iom.queue.Publish(topic, string(d))
}

func (iom *iomanager) loop() {
	log.Info("Starting IO Manager Loop")

	lastHC := time.Now()

	for iom.running {
		iom.devLock.Lock()
		if time.Since(lastHC) > time.Second {
			_ = iom.dev.HealthCheck()
			lastHC = time.Now()
			iom.status, _ = iom.dev.Status()
			iom.updateStatus()
		}
		lines := iom.dev.ReadAllLines()
		for _, v := range lines {
			// skip HC
			if !strings.Contains(v, "Health Check") {
				log.Info("[SERIAL] %s", v)
			}
		}
		iom.devLock.Unlock()
		time.Sleep(time.Millisecond * 100)
	}

	log.Info("Stopping IO Manager Loop")
}

func (iom *iomanager) Close() error {
	iom.running = false
	_ = iom.queue.Close()
	return iom.dev.Close()
}
