package config

type ServerOpts struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type MysqlOpts struct {
	Address         string `mapstructure:"address" yaml:"address"`
	UserName        string `mapstructure:"username" yaml:"username"`
	Password        string `mapstructure:"password" yaml:"password"`
	DBName          string `mapstructure:"dbname" yaml:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}
