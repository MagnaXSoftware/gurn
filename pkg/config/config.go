package config

type Config struct {
	BindAddr string `hcl:"bind_addr"`
	Hostname string `hcl:"hostname,optional"`

	Database string `hcl:"database,optional"`
}

func DefaultConfig() *Config {
	conf := &Config{
		Database: "gurn.db",
	}
	return conf
}
