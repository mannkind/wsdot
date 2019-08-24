package main

import (
	"fmt"
	"reflect"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

type mqttClient struct {
	twomqtt.Observer
	*twomqtt.MQTTProxy
	mqttClientConfig
}

func newMQTTClient(mqttClientCfg mqttClientConfig, client *twomqtt.MQTTProxy) *mqttClient {
	c := mqttClient{
		MQTTProxy:        client,
		mqttClientConfig: mqttClientCfg,
	}

	c.Initialize(
		c.onConnect,
		c.onDisconnect,
	)

	c.LogSettings()

	return &c
}

func (c *mqttClient) run() {
	c.Run()
}

func (c *mqttClient) onConnect(client mqtt.Client) {
	log.Info("Connected to MQTT")
	c.Publish(c.AvailabilityTopic(), "online")
	c.publishDiscovery()
}

func (c *mqttClient) onDisconnect(client mqtt.Client, err error) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error("Disconnected from MQTT")
}

func (c *mqttClient) publishDiscovery() {
	if !c.Discovery {
		return
	}

	log.Info("MQTT discovery publishing")

	for _, travelTimeSlug := range c.TravelTimeMapping {
		sensor := strings.ToLower(travelTimeSlug)
		mqd := c.NewMQTTDiscovery("", sensor, "sensor")
		mqd.Icon = "mdi:car"
		mqd.UnitOfMeasurement = "min"

		c.PublishDiscovery(mqd)
	}

	log.Info("Finished MQTT discovery publishing")
}

func (c *mqttClient) ReceiveCommand(cmd twomqtt.Command, e twomqtt.Event) {}
func (c *mqttClient) ReceiveState(e twomqtt.Event) {
	if e.Type != reflect.TypeOf(wsdotTravelTime{}) {
		msg := "Unexpected event type; skipping"
		log.WithFields(log.Fields{
			"type": e.Type,
		}).Error(msg)
		return
	}

	info := e.Payload.(wsdotTravelTime)
	travelTimeID := fmt.Sprintf("%d", info.TravelTimeID)
	travelTimeSlug := c.TravelTimeMapping[travelTimeID]

	topic := c.StateTopic("", travelTimeSlug)
	payload := fmt.Sprintf("%d", info.CurrentTime)

	if info.Distance == 0 {
		payload = "Closed"
	} else if info.CurrentTime == 0 {
		payload = "Unknown"
	}

	c.Publish(topic, payload)
}
