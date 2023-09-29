package usecases

import (
	"fmt"
	"os"
	"theo303/neon-pricer/internal/svg"
)

func GetPerimeter(filepath string) (float64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, fmt.Errorf("opening file: %w", err)
	}

	forms, err := svg.RetrieveForms(file)
	if err != nil {
		return 0, fmt.Errorf("retrieving forms from svg file: %w", err)
	}
	fmt.Printf("%+v\n", forms)

	var totalPerimeter float64
	for _, form := range forms {
		totalPerimeter += form.GetPerimeter()
	}

	return totalPerimeter, nil
}
