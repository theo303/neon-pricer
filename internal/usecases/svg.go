package usecases

import (
	"fmt"
	"os"

	"theo303/neon-pricer/internal/domain"
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

func GetBounds(formsGroups map[string][]svg.Form) (map[string]svg.Bounds, error) {
	bounds := make(map[string]svg.Bounds)
	for id, forms := range formsGroups {
		var b svg.Bounds
		for i, form := range forms {
			nb, err := form.Bounds()
			if err != nil {
				return nil, fmt.Errorf("retrieving bounds on form n %d: %w", i, err)
			}
			if i == 0 {
				b = nb
			} else {
				b = b.Expand(nb)
			}
		}
		bounds[id] = b
	}
	return bounds, nil
}

// GetSizes retrieves lengths, height and width for each group of form and scales them.
func GetSizes(formsGroups map[string][]svg.Form, scale float64) (map[string]domain.Size, error) {
	lengths, err := GetLengths(formsGroups)
	if err != nil {
		return nil, fmt.Errorf("computing lengths: %w", err)
	}
	bounds, err := GetBounds(formsGroups)
	if err != nil {
		return nil, fmt.Errorf("computing bounds: %w", err)
	}
	sizes := make(map[string]domain.Size)
	for id := range formsGroups {
		length, ok := lengths[id]
		if !ok {
			return nil, fmt.Errorf("missing id %s in lengths map", id)
		}
		bound, ok := bounds[id]
		if !ok {
			return nil, fmt.Errorf("missing id %s in bounds map", id)
		}
		sizes[id] = domain.Size{
			Length:   length * 1000 / scale,
			LengthPx: length,
			Height:   bound.Height() * 1000 / scale,
			Width:    bound.Width() * 1000 / scale,
		}
	}
	return sizes, nil
}
