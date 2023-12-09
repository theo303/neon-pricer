package configuration

import "github.com/spf13/viper"

const configPath = "./"

type Configuration struct {
	Scale float64 `mapstructure:"scale"`
}

// Load reads configuration from file.
func Load() (Configuration, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	config := Configuration{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}
