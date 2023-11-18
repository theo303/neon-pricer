package usecases

import (
	"fmt"
	"os"
	"theo303/neon-pricer/internal/svg"
)

// ParseSVGFile parses an svg file and returns a map containing forms found in each groups.
func ParseSVGFile(filepath string, groupID string) (map[string][]svg.Form, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	formsGroups, err := svg.RetrieveForms(file, groupID)
	if err != nil {
		return nil, fmt.Errorf("retrieving forms from svg file: %w", err)
	}

	return formsGroups, nil
}

func GetLengths(formsGroups map[string][]svg.Form) (map[string]float64, error) {
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

func GetSizes(formsGroups map[string][]svg.Form) (map[string]float64, error) {
	return nil, nil
}
