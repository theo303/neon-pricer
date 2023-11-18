package usecases

import (
	"fmt"
	"os"
	"theo303/neon-pricer/internal/svg"
)

func GetLengths(filepath string, groupID string) (map[string]float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	formsGroups, err := svg.RetrieveForms(file, groupID)
	if err != nil {
		return nil, fmt.Errorf("retrieving forms from svg file: %w", err)
	}

	lengths := make(map[string]float64)
	for id, forms := range formsGroups {
		for _, form := range forms {
			l, err := form.Length()
			if err != nil {
				return nil, err
			}
			lengths[id] += l
		}
	}

	return lengths, nil
}
