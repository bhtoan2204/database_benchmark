package config

type ElasticSearchConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Config struct {
	ElasticSearchConfig ElasticSearchConfig `mapstructure:"elasticsearch"`
}
