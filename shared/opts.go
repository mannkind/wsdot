package shared

import (
	"reflect"

	"github.com/caarlos0/env/v6"
	"github.com/mannkind/twomqtt"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	TravelTimeMapping map[string]string `env:"WSDOT_TRAVELTIMEMAPPING" envDefault:"132:seattle2everett,31:seattle2renton"`
}

// NewOpts creates a Opts based on environment variables
func NewOpts() Opts {
	c := Opts{}

	// Manually parse the address:name mapping
	if err := env.ParseWithFuncs(&c, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(map[string]string{}): twomqtt.SimpleKVMapParser(":", ","),
	}); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	return c
}
