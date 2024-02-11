package usecases

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/domain"
)

const defaultSiliconeSize = 12

var mmRegexp = regexp.MustCompile(`(?i)(\d+)MM`)

type LayerPrice struct {
	SiliconePrice float64
	LEDPrice      float64
	PlexiPrice    float64
}

type Price map[string]LayerPrice

func GetPrice(config conf.Pricing, sizes map[string]domain.Size, plexi string) (Price, error) {
	price := make(map[string]LayerPrice)
	for id, size := range sizes {
		id = strings.ToUpper(id)
		if id == "DECOUPE" {
			area := size.Height / 1000 * size.Width / 1000
			price[id] = LayerPrice{
				PlexiPrice: domain.Round(getPlexiPricing(config.Plexis, plexi) * area),
			}
		}

		siliconeSize, err := getSiliconeSize(id)
		if err != nil {
			return Price{}, fmt.Errorf("retrieving silicone size: %w", err)
		}
		if siliconeSize == 0 {
			continue
		}

		siliconePrice, err := getSiliconePricing(config.Silicones, siliconeSize)
		if err != nil {
			return Price{}, fmt.Errorf("retrieving silicone price: %w", err)
		}
		price[id] = LayerPrice{
			SiliconePrice: domain.Round(siliconePrice * size.Length / 1000),
			LEDPrice:      domain.Round(getLEDPricing(config.LEDs, id) * size.Length / 1000),
		}
	}
	return price, nil
}

func getSiliconeSize(id string) (int, error) {
	if slices.Contains([]string{"RGB", "PIXEL"}, id) {
		return defaultSiliconeSize, nil
	}
	if submatchs := mmRegexp.FindStringSubmatch(id); len(submatchs) > 1 {
		size, err := strconv.Atoi(submatchs[1])
		if err != nil {
			return 0, fmt.Errorf("%s could not be converted to int: %w", submatchs[1], err)
		}
		return size, nil
	}
	return 0, nil
}

func getSiliconePricing(pricings []conf.Silicone, size int) (float64, error) {
	for _, pricingSilicone := range pricings {
		if pricingSilicone.SizeMm == size {
			return pricingSilicone.PricePerMeter, nil
		}
	}
	return 0, fmt.Errorf("no pricing could be found for size %dMM", size)
}

func getLEDPricing(pricings []conf.LED, id string) float64 {
	var defaultPricing float64
	for _, pricingLed := range pricings {
		switch pricingLed.Name {
		case id:
			return pricingLed.PricePerMeter
		case "couleur":
			defaultPricing = pricingLed.PricePerMeter
		}
	}
	return defaultPricing
}

func getPlexiPricing(pricings []conf.Plexi, name string) float64 {
	var defaultPricing float64
	fmt.Println(name)
	for _, pricingPlexi := range pricings {
		switch pricingPlexi.Name {
		case name:
			return pricingPlexi.PricePerMeterSquare
		case "incolore":
			defaultPricing = pricingPlexi.PricePerMeterSquare
		}
	}
	return defaultPricing
}
