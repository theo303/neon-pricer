package usecases

import (
	"fmt"
	"os"
	"theo303/neon-pricer/internal/svg"
)

func GetLength(filepath string, groupID string) (float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("opening file: %w", err)
	}

	forms, err := svg.RetrieveForms(file, groupID)
	if err != nil {
		return 0, fmt.Errorf("retrieving forms from svg file: %w", err)
	}

	var totalPerimeter float64
	for _, form := range forms {
		p, err := form.Length()
		if err != nil {
			return 0, err
		}
		totalPerimeter += p
	}

	return totalPerimeter, nil
}
