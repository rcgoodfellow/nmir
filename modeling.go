package nmir

import (
	"fmt"
)

func Unit(qual string, value float64, unit string) map[string]interface{} {
	return map[string]interface{}{
		fmt.Sprintf("value%s", qual): value,
		"unit": unit,
	}
}
