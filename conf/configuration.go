package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

const configPath = "./"

type Silicone struct {
	SizeMm        int     `mapstructure:"size"`
	PricePerMeter float64 `mapstructure:"price"`
}

type LED struct {
	Name          string  `mapstructure:"name"`
	PricePerMeter float64 `mapstructure:"price"`
}

type Plexi struct {
	Name                string  `mapstructure:"name"`
	PricePerMeterSquare float64 `mapstructure:"price"`
}

type Controler struct {
	Name  string  `mapstructure:"name"`
	Price float64 `mapstructure:"price"`
}

type PowerSupply struct {
	Amp   string  `mapstructure:"amp"`
	Price float64 `mapstructure:"price"`
}

type Pricing struct {
	Silicones     []Silicone    `mapstructure:"silicones"`
	LEDs          []LED         `mapstructure:"leds"`
	Plexis        []Plexi       `mapstructure:"plexis"`
	Controlers    []Controler   `mapstructure:"controlers"`
	PowerSupplies []PowerSupply `mapstructure:"power_supplies"`
}

type Configuration struct {
	Pricing `mapstructure:",squash"`
	Scale   float64 `mapstructure:"scale"`
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

	fmt.Printf("configuration loaded: %+v\n", config)
	return config, nil
}
