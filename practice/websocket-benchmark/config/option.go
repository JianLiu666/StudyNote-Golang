package config

type ServerOpts struct {
	Addr string `mapstructure:"addr" yaml:"addr"`
	Port string `mapstructure:"port" yaml:"port"`
}

type SimulationOpts struct {
	NumClients  int `mapstructure:"num_clients" yaml:"num_clients"`
	NumMessages int `mapstructure:"num_messages" yaml:"num_messages"`
}
