package config

type ServerOpts struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type RedisOpts struct {
	Address            string `mapstructure:"address"`
	Password           string `mapstructure:"password"`
	DB                 int    `mapstructure:"db"`
	PoolSize           int    `mapstructure:"pool_size"`
	ListCollectionName string `mapstructure:"list_collection_name"`
	PageExpirationTime int    `mapstructure:"page_expiration_time"`
	PageKeyPrefix      string `mapstructure:"page_key_prefix"`
}
