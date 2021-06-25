package mpstmetadatataggingprocessor

import (
	"errors"
	"fmt"
	"github.com/scylladb/go-set/strset"
	"go.opentelemetry.io/collector/config"
)

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	Roles                    map[string]metadataTag `mapstructure:"roles"`
}

func (c *Config) Validate() error {
	if len(c.Roles) == 0 {
		return errors.New("no roles defined in config")
	}
	definedIdentifiers := strset.New()
	for roleName, roleData := range c.Roles {
		roleIdentifier := roleData.Name
		if definedIdentifiers.Has(roleIdentifier) {
			return fmt.Errorf("duplicate role identifier: %s", roleName)
		}
		definedIdentifiers.Add(roleIdentifier)
	}
	return nil
}

type metadataTag struct {
	Name string `mapstructure:"name"`
}
