package main

import (
	"github.com/quan-to/slog"
	"github.com/racerxdl/esp32-shift/internal/espshift"
	"github.com/racerxdl/esp32-shift/pkg/proto"
	"time"
)

var log = slog.Scope("Test Program")

func main() {
	var status []shift.IOLevel
	esp, err := espshift.MakeESPShift("/dev/ttyUSB0")

	if err != nil {
		panic(err)
	}

	defer esp.Close()

	lastData := 1

	for err == nil {
		err = esp.HealthCheck()
		if err != nil {
			break
		}

		status, err = esp.Status()
		if err != nil {
			break
		}

		log.Info("Board Status: %s", status)

		err = esp.SetByte(0, uint8(lastData&0xFF))
		if err != nil {
			break
		}
		err = esp.SetByte(1, uint8(lastData>>8))
		if err != nil {
			break
		}

		lastData <<= 1
		if lastData > 65535 {
			lastData = 1
		}
		time.Sleep(time.Millisecond * 10)

		d := esp.ReadAllLines()

		for _, v := range d {
			log.Info("[SERIAL] %s", v)
		}
	}

	if err != nil {
		panic(err)
	}
}
