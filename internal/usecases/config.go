package usecases

import (
	"fmt"
	"strconv"
	"strings"
	"theo303/neon-pricer/conf"
)

const (
	scaleParam       = "scale"
	siliconeParam    = "silic"
	ledParam         = "led"
	plexiParam       = "plexi"
	controlerParam   = "controler"
	powerSupplyParam = "powersupply"
)

func UpdateConfigWithPostForm(config *conf.Configuration, body []byte) (*conf.Configuration, error) {
	pairs := strings.Split(string(body), "&")
	values := make(map[string]float64)
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("pair %s does not contain 2 parts", pair)
		}
		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", parts[1], err)
		}
		values[parts[0]] = value
	}

	config.Scale = values[scaleParam]
	for idx, s := range config.Silicones {
		config.Silicones[idx].PricePerMeter = values[fmt.Sprintf("%s-%d", siliconeParam, s.SizeMm)]
	}
	for idx, l := range config.LEDs {
		config.LEDs[idx].PricePerMeter = values[fmt.Sprintf("%s-%s", ledParam, l.Name)]
	}
	for idx, p := range config.Plexis {
		config.Plexis[idx].PricePerMeterSquare = values[fmt.Sprintf("%s-%s", plexiParam, p.Name)]
	}
	for idx, c := range config.Controlers {
		config.Controlers[idx].Price = values[fmt.Sprintf("%s-%s", controlerParam, c.Name)]
	}
	for idx, ps := range config.PowerSupplies {
		config.PowerSupplies[idx].Price = values[fmt.Sprintf("%s-%s", powerSupplyParam, ps.Amp)]
	}
	return config, nil
}
