package api

import "testing"

func TestSanitizeImageSize(t *testing.T) {
	inputs := []uint64{
		536870912,  // 0.5 GB
		1073741824, // 1   GB
		1610612736, // 1.5 GB
	}
	outputs := []int{
		1,
		1,
		2,
	}

	for i, input := range inputs {
		o := sanitizeImageSize(input)
		if outputs[i] != o {
			t.Errorf("failed to sanitize (input: %d)", input)
		}
	}
}
