package queue

import "github.com/eclipse/paho.mqtt.golang"

type OnMessage func(mqtt.Message)

type Manager interface {
	SetOnMessage(cb OnMessage)
	Publish(topic string, payload interface{}) error
}
