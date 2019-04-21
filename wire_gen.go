// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/mannkind/paho.mqtt.golang.ext/cfg"
	"github.com/mannkind/paho.mqtt.golang.ext/di"
)

// Injectors from wire.go:

func initialize() *bridge {
	mqttConfig := cfg.NewMQTTConfig()
	mainConfig := newConfig(mqttConfig)
	mqttFuncWrapper := di.NewMQTTFuncWrapper()
	mainMqttClient := newMQTTClient(mainConfig, mqttFuncWrapper)
	mainClient := newClient(mainConfig)
	mainBridge := newBridge(mainConfig, mainMqttClient, mainClient)
	return mainBridge
}
