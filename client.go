package main

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

const (
	travelTimeURL = "https://www.wsdot.wa.gov/Traffic/api/TravelTimes/TravelTimesREST.svc/GetTravelTimeAsJson"
)

type client struct {
	observers map[observer]struct{}

	lookupInterval time.Duration
	travelTimes    map[string]string
	secret         string
}

func newClient(config *config) *client {
	c := client{
		observers: map[observer]struct{}{},

		lookupInterval: config.LookupInterval,
		secret:         config.Secret,

		travelTimes: map[string]string{},
	}

	// Create a mapping between travel time ids and names
	for _, m := range config.TravelTimeMapping {
		parts := strings.Split(m, ":")
		if len(parts) != 2 {
			continue
		}

		travelTimeID := parts[0]
		travelTimeName := parts[1]
		c.travelTimes[travelTimeID] = travelTimeName
	}

	return &c
}

func (c *client) run() {
	go c.loop(false)
}

func (c *client) register(l observer) {
	c.observers[l] = struct{}{}
}

func (c *client) publish(e event) {
	for o := range c.observers {
		o.receiveState(e)
	}
}

func (c *client) loop(once bool) {
	for {
		log.Info("Beginning lookup")
		for travelTimeID, travelTimeSlug := range c.travelTimes {
			if info, err := c.lookup(travelTimeID); err == nil {
				c.publish(event{
					version: 1,
					key:     travelTimeSlug,
					data:    c.adapt(info),
				})
			}
		}
		log.Info("Ending lookup")

		if once {
			break
		}

		time.Sleep(c.lookupInterval)
	}
}

func (c *client) lookup(travelTimeID string) (*wsdotTravelTime, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&wsdotTravelTime{}).
		SetQueryParams(map[string]string{
			"AccessCode":   c.secret,
			"TravelTimeID": travelTimeID,
		}).
		Get(travelTimeURL)

	if err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"travelTimeID": travelTimeID,
		}).Error("Unable to lokup the travel time specified")
		return nil, err
	}

	return resp.Result().(*wsdotTravelTime), nil
}

func (c *client) adapt(info *wsdotTravelTime) eventData {
	return eventData{
		CurrentTime:  info.CurrentTime,
		Distance:     info.Distance,
		TravelTimeID: info.TravelTimeID,
	}
}
