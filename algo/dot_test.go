package algo

import (
	"testing"
)

func TestDotBuilderEqualDistribution(t *testing.T) {
	testCases := []struct {
		Name string
		Size int
	}{
		{
			Name: "1",
			Size: 1,
		},
		{
			Name: "2",
			Size: 2,
		},
		{
			Name: "3",
			Size: 3,
		},
		{
			Name: "4",
			Size: 4,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			maxValue := byte(testCase.Size * testCase.Size)

			values := newDotBuilder(testCase.Size, false).generateMatrix(0, maxValue-1)

			expectedValues := make(map[byte]struct{})
			for i := byte(0); i < maxValue; i++ {
				expectedValues[i] = struct{}{}
			}

			for x := 0; x < testCase.Size; x++ {
				for y := 0; y < testCase.Size; y++ {
					delete(expectedValues, values.Get(x, y))
				}
			}

			if len(expectedValues) != 0 {
				t.Fail()
			}
		})
	}
}
