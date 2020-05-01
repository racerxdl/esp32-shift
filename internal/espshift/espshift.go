package espshift

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jacobsa/go-serial/serial"
	"github.com/quan-to/slog"
	"github.com/racerxdl/esp32-shift/pkg/proto"
	"io"
	"strconv"
	"strings"
	"time"
)

var log = slog.Scope("SerialDevice")

const (
	bsHeader = "( OK) BS["
)

type serialDevice struct {
	port       io.ReadWriteCloser
	lineReader *bufio.Reader
}

func MakeESPShift(device string) (shift.ESPShift, error) {
	options := serial.OpenOptions{
		PortName:              device,
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 50,
	}

	log.Debug("Opening port %s", device)
	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	esp := &serialDevice{
		port:       port,
		lineReader: bufio.NewReader(port),
	}

	log.Debug("Checking liveness")
	if !esp.checkStatus() {
		return nil, fmt.Errorf("error checking device liveness")
	}

	log.Debug("Done!")
	return esp, nil
}

func (sd *serialDevice) checkStatus() bool {
	var err error
	var line []byte

	// Read all lines in buffer
	log.Debug("Reading lines")
	d := sd.ReadAllLines()

	for _, v := range d {
		log.Info("[SERIAL] %s", v)
	}

	// Issue a HealthCheck Command and wait for return
	err = sd.HealthCheck()
	if err != nil {
		log.Error("Error sending HealthCheck: %s", err)
		return false
	}
	time.Sleep(time.Millisecond * 10)

	err = nil
	for err == nil {
		line, _, err = sd.lineReader.ReadLine()
		if err != nil {
			log.Error("Error getting HealthCheck line: %s", err)
			return false
		}
		sline := string(line)
		if sline == "( OK) Health Check OK" {
			return true
		}
	}

	return false
}

func (sd *serialDevice) issueCommand(msg *shift.CmdMsg) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	txData := make([]byte, len(data)+2)
	binary.LittleEndian.PutUint16(txData[:2], uint16(len(data)))
	copy(txData[2:], data)

	_, err = sd.port.Write(txData)
	return err
}

func (sd *serialDevice) Reset() error {
	log.Debug("Sending RESET")
	cmd := &shift.CmdMsg{
		Cmd: shift.CmdMsg_Reset,
	}

	return sd.issueCommand(cmd)
}

func (sd *serialDevice) HealthCheck() error {
	//log.Debug("Sending HEALTHCHECK")
	cmd := &shift.CmdMsg{
		Cmd: shift.CmdMsg_HealthCheck,
	}

	return sd.issueCommand(cmd)
}

func (sd *serialDevice) SetPin(pin uint8, val shift.IOLevel) error {
	log.Debug("Sending SETPIN(%d, %s)", pin, val)
	cmd := &shift.CmdMsg{
		Cmd:  shift.CmdMsg_SetPin,
		Data: []byte{pin, uint8(val)},
	}

	return sd.issueCommand(cmd)
}

func (sd *serialDevice) SetByte(byteNum uint8, val uint8) error {
	log.Debug("Sending SETBYTE(%d, %08b)", byteNum, val)
	cmd := &shift.CmdMsg{
		Cmd:  shift.CmdMsg_SetByte,
		Data: []byte{byteNum, val},
	}

	return sd.issueCommand(cmd)
}

func (sd *serialDevice) Status() (shift.BoardStatus, error) {
	//log.Debug("Sending STATUS")
	var line []byte
	cmd := &shift.CmdMsg{
		Cmd: shift.CmdMsg_Status,
	}

	err := sd.issueCommand(cmd)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 10)

	err = nil
	bsLine := ""
	for err == nil {
		line, _, err = sd.lineReader.ReadLine()
		if err != nil {
			log.Error("Error getting Status Line: %s", err)
			return nil, err
		}
		sline := string(line)
		if len(sline) > len(bsHeader) && sline[:len(bsHeader)] == bsHeader {
			bsLine = sline
			break
		}
	}

	bsLine = bsLine[len(bsHeader):] // Remove bsHeader
	bsLine = bsLine[:len(bsLine)-1] // Remove last ]

	vals := strings.Split(bsLine, ",")
	levels := make([]shift.IOLevel, len(vals))

	for i, v := range vals {
		vi, _ := strconv.ParseInt(strings.Trim(v, " \r\n"), 10, 32)
		levels[i] = shift.IOLevel(vi & 0xFF)
	}

	return levels, nil
}

func (sd *serialDevice) Close() error {
	return sd.port.Close()
}

func (sd *serialDevice) ReadAllLines() []string {
	var err error
	var lines []string
	var line []byte

	err = nil
	for err == nil {
		line, _, err = sd.lineReader.ReadLine()
		if err != nil {
			break
		}
		sline := string(line)
		lines = append(lines, sline)
	}

	return lines
}
