package queue

import (
	"github.com/eclipse/paho.mqtt.golang"
	"io"
)

type OnMessage func(mqtt.Message)

type Manager interface {
	io.Closer
	SetOnMessage(cb OnMessage)
	Subscribe(topic string) error
	Publish(topic string, payload interface{}) error
}
