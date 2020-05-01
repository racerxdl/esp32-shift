package main

import (
	"github.com/quan-to/slog"
	"github.com/racerxdl/esp32-shift/internal/espshift"
	"github.com/racerxdl/esp32-shift/internal/ioman"
	"github.com/racerxdl/esp32-shift/internal/iqueue"
	"github.com/racerxdl/esp32-shift/pkg/queue"
)

var log = slog.Scope("QueueShift")

func main() {

	LoadConfig()

	q, err := iqueue.MakeQueueManager(queue.MQTTConfig{
		MQTTUsername: config.MQTTUsername,
		MQTTPassword: config.MQTTPassword,
		MQTTServer:   config.MQTTServer,
		CloseQueue:   config.CloseQueue,
	})

	if err != nil {
		log.Fatal(err)
	}

	dev, err := espshift.MakeESPShift(config.SerialPort)

	if err != nil {
		log.Fatal(err)
	}

	io := ioman.MakeIOManager(dev, config.BaseTopic, q)

	select {}

	_ = io.Close()
}
