package main

import (
	"strconv"
	"time"

	"github.com/nDenerserve/SmartPi/src/smartpi"
	MQTT "github.com/rvolz/gomosquittogo/core"
)

func newMQTTClient(c *smartpi.Config) (mqttclient MQTT.Client) {

	mqttclient := NewClient(c.MQTTbroker, nil)
	defer mqttclient.Close()

	mqttclient.Connect()
	//mqttclient.SendString("my topic","this is my message")

	return mqttclient
}

func publishMQTTReadouts(c *smartpi.Config, mqttclient MQTT.Client, values *smartpi.ADE7878Readout) {

	if mqttclient.clientConnected {

		mqttclient.SendString(c.MQTTtopic+"/I4", values.Current[smartpi.PhaseN])
		for _, p := range smartpi.MainPhases {
			label := p.PhaseNumber()
			mqttclient.SendString(c.MQTTtopic+"/I"+label, values.Current[p])
			mqttclient.SendString(c.MQTTtopic+"/V"+label, values.Voltage[p])
			mqttclient.SendString(c.MQTTtopic+"/P"+label, values.ActiveWatts[p])
			mqttclient.SendString(c.MQTTtopic+"/COS"+label, values.CosPhi[p])
			mqttclient.SendString(c.MQTTtopic+"/F"+label, values.Frequency[p])
		}
	}
}
