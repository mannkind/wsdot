package source

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/mannkind/wsdot2mqtt/shared"
	log "github.com/sirupsen/logrus"
)

// Opts is for package related settings
type Opts struct {
	shared.Opts
	Secret         string        `env:"WSDOT_SECRET,required"`
	LookupInterval time.Duration `env:"WSDOT_LOOKUPINTERVAL"    envDefault:"3m"`
}

// NewOpts creates a Opts based on environment variables
func NewOpts(opts shared.Opts) Opts {
	c := Opts{
		Opts: opts,
	}

	if err := env.Parse(&c); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to unmarshal configuration")
	}

	return c
}
