package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

func ParseConfig(src []byte) (*Config, error) {
	conf := DefaultConfig()

	err := hclsimple.Decode("config.hcl", src, nil, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
